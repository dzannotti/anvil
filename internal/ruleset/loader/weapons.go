package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"anvil/internal/expression"
	"anvil/internal/ruleset/basic"
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

func LoadWeapons(dataDir string) (map[string]func() *basic.Weapon, error) {
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

		if !strings.HasPrefix(archetype, "_") {
			weaponDefs[archetype] = weaponData
		}
	}

	return weaponDefs, nil
}

func createWeaponFactories(weaponDefs map[string]WeaponData) map[string]func() *basic.Weapon {
	weaponFactories := make(map[string]func() *basic.Weapon)

	for archetype, weaponData := range weaponDefs {
		data := weaponData
		arch := archetype

		weaponFactories[arch] = func() *basic.Weapon {
			return createWeapon(arch, data)
		}
	}

	return weaponFactories
}

func createWeapon(archetype string, data WeaponData) *basic.Weapon {
	damageExpr := expression.Expression{Rng: expression.DefaultRoller{}}
	for _, dmg := range data.Damage {
		if err := parseDamageFormula(dmg.Formula, data.Name, dmg.Kind, &damageExpr); err != nil {
			panic(fmt.Sprintf("invalid damage formula '%s' for weapon '%s': %v", dmg.Formula, archetype, err))
		}
	}

	weaponTags := make([]tag.Tag, len(data.WeaponTags))
	for i, tagStr := range data.WeaponTags {
		weaponTags[i] = tag.FromString(tagStr)
	}

	return basic.NewWeapon(
		archetype,
		uuid.New().String(),
		data.Name,
		damageExpr,
		tag.NewContainer(weaponTags...),
		data.Reach,
	)
}

func parseDamageFormula(formula, weaponName, kind string, expr *expression.Expression) error {
	diceRe := regexp.MustCompile(`^(\d+)d(\d+)$`)
	if matches := diceRe.FindStringSubmatch(formula); len(matches) == 3 {
		times, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("invalid number of dice: %s", matches[1])
		}

		sides, err := strconv.Atoi(matches[2])
		if err != nil {
			return fmt.Errorf("invalid number of sides: %s", matches[2])
		}

		expr.AddDamageDice(times, sides, weaponName, tag.NewContainer(tag.FromString(kind)))
		return nil
	}

	if constant, err := strconv.Atoi(formula); err == nil {
		expr.AddDamageConstant(constant, weaponName, tag.NewContainer(tag.FromString(kind)))
		return nil
	}

	return fmt.Errorf("invalid damage formula format: %s (expected format like '1d4', '2d6', or '5')", formula)
}
