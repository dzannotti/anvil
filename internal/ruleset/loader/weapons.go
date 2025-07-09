package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"anvil/internal/core/tags"
	weapon "anvil/internal/ruleset/items/weapons"
	"anvil/internal/tag"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type DamageData struct {
	Formula string `yaml:"formula"`
	Kind    string `yaml:"kind"`
}

type WeaponData struct {
	Name       string       `yaml:"name"`
	Damage     []DamageData `yaml:"damage"`
	WeaponTags []string     `yaml:"weapon_tags"`
	Reach      int          `yaml:"reach"`
}

type WeaponsFile struct {
	Weapons map[string]WeaponData `yaml:"weapons"`
}

func LoadWeapons(dataDir string) (map[string]func() *weapon.Weapon, error) {
	weaponsDir := filepath.Join(dataDir, "ruleset", "weapons")

	weaponDefs, err := loadWeaponFiles(weaponsDir)
	if err != nil {
		return nil, err
	}

	return createWeaponFactories(weaponDefs), nil
}

func loadWeaponFiles(weaponsDir string) (map[string]WeaponData, error) {
	pattern := filepath.Join(weaponsDir, "*.yml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob weapons files: %w", err)
	}

	weaponDefs := make(map[string]WeaponData)

	for _, match := range matches {
		fileData, readErr := os.ReadFile(match)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read weapon file %s: %w", match, readErr)
		}

		var weaponData WeaponData
		if unmarshalErr := yaml.Unmarshal(fileData, &weaponData); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to parse weapon YAML %s: %w", match, unmarshalErr)
		}

		basename := filepath.Base(match)
		archetype := strings.TrimSuffix(basename, ".yml")

		// Skip base archetype files (starting with _)
		if !strings.HasPrefix(archetype, "_") {
			weaponDefs[archetype] = weaponData
		}
	}

	return weaponDefs, nil
}

func createWeaponFactories(weaponDefs map[string]WeaponData) map[string]func() *weapon.Weapon {
	weaponFactories := make(map[string]func() *weapon.Weapon)

	for archetype, weaponData := range weaponDefs {
		data := weaponData
		arch := archetype

		weaponFactories[arch] = func() *weapon.Weapon {
			return createWeapon(arch, data)
		}
	}

	return weaponFactories
}

func createWeapon(archetype string, data WeaponData) *weapon.Weapon {
	damageEntries := make([]weapon.DamageEntry, 0, len(data.Damage))
	for _, dmg := range data.Damage {
		times, sides, formulaErr := parseDiceFormula(dmg.Formula)
		if formulaErr != nil {
			panic(fmt.Sprintf("invalid dice formula '%s' for weapon '%s': %v", dmg.Formula, archetype, formulaErr))
		}
		damageEntries = append(damageEntries, weapon.DamageEntry{
			Times: times,
			Sides: sides,
			Kind:  stringToTag(dmg.Kind),
		})
	}

	weaponTags := make([]tag.Tag, len(data.WeaponTags))
	for i, tagStr := range data.WeaponTags {
		weaponTags[i] = stringToTag(tagStr)
	}

	return weapon.NewWeapon(
		archetype,
		uuid.New().String(),
		data.Name,
		damageEntries,
		tag.NewContainer(weaponTags...),
		data.Reach,
	)
}

func parseDiceFormula(formula string) (int, int, error) {
	// Parse formulas like "1d4", "2d6", etc.
	re := regexp.MustCompile(`^(\d+)d(\d+)$`)
	matches := re.FindStringSubmatch(formula)
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid dice formula format, expected format like '1d4' or '2d6'")
	}

	times, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid number of dice: %s", matches[1])
	}

	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid number of sides: %s", matches[2])
	}

	return times, sides, nil
}

var knownTags = map[string]tag.Tag{
	"Damage.Kind.Piercing":    tags.Piercing,
	"Damage.Kind.Slashing":    tags.Slashing,
	"Damage.Kind.Bludgeoning": tags.Bludgeoning,
	"Damage.Kind.Fire":        tags.Fire,
	"Damage.Kind.Poison":      tags.Poison,
	"Damage.Kind.Radiant":     tags.Radiant,
	"Melee":                   tags.Melee,
	"Ranged":                  tags.Ranged,
	"Item.Weapon.Simple":      tags.SimpleWeapon,
	"Item.Weapon.Martial":     tags.MartialWeapon,
	"Item.Weapon.Martial.Axe": tags.MartialAxe,
	"Item.Weapon.Natural":     tags.NaturalWeapon,
	"Item.Weapon.Finesse":     tags.Finesse,
}

func stringToTag(tagStr string) tag.Tag {
	if knownTag, exists := knownTags[tagStr]; exists {
		return knownTag
	}
	return tag.FromString(tagStr)
}
