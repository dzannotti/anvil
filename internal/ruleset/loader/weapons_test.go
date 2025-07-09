package loader

import (
	"os"
	"testing"

	"anvil/internal/expression"
	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDamageFormula_DiceFormulas(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		weaponName  string
		kind        string
		expectTimes int
		expectSides int
		expectError bool
	}{
		{
			name:        "1d4",
			formula:     "1d4",
			weaponName:  "Dagger",
			kind:        "Damage.Kind.Piercing",
			expectTimes: 1,
			expectSides: 4,
			expectError: false,
		},
		{
			name:        "2d6",
			formula:     "2d6",
			weaponName:  "Greataxe",
			kind:        "Damage.Kind.Slashing",
			expectTimes: 2,
			expectSides: 6,
			expectError: false,
		},
		{
			name:        "10d10",
			formula:     "10d10",
			weaponName:  "Fireball",
			kind:        "Damage.Kind.Fire",
			expectTimes: 10,
			expectSides: 10,
			expectError: false,
		},
		{
			name:        "invalid dice format",
			formula:     "1d",
			weaponName:  "Invalid",
			kind:        "Damage.Kind.Piercing",
			expectError: true,
		},
		{
			name:        "invalid dice format 2",
			formula:     "d6",
			weaponName:  "Invalid",
			kind:        "Damage.Kind.Piercing",
			expectError: true,
		},
		{
			name:        "invalid dice format 3",
			formula:     "1d6d8",
			weaponName:  "Invalid",
			kind:        "Damage.Kind.Piercing",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := expression.Expression{Rng: expression.DefaultRoller{}}
			err := parseDamageFormula(tt.formula, tt.weaponName, tt.kind, &expr)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, expr.Components, 1)

			component := expr.Components[0]
			assert.Equal(t, tt.expectTimes, component.Times)
			assert.Equal(t, tt.expectSides, component.Sides)
			assert.Equal(t, tt.weaponName, component.Source)
			assert.True(t, component.Tags.HasTag(tag.FromString(tt.kind)))
		})
	}
}

func TestParseDamageFormula_Constants(t *testing.T) {
	tests := []struct {
		name        string
		formula     string
		weaponName  string
		kind        string
		expectValue int
		expectError bool
	}{
		{
			name:        "single digit",
			formula:     "5",
			weaponName:  "Magic Weapon",
			kind:        "Damage.Kind.Force",
			expectValue: 5,
			expectError: false,
		},
		{
			name:        "double digit",
			formula:     "15",
			weaponName:  "Big Magic",
			kind:        "Damage.Kind.Fire",
			expectValue: 15,
			expectError: false,
		},
		{
			name:        "zero",
			formula:     "0",
			weaponName:  "Harmless",
			kind:        "Damage.Kind.Bludgeoning",
			expectValue: 0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := expression.Expression{Rng: expression.DefaultRoller{}}
			err := parseDamageFormula(tt.formula, tt.weaponName, tt.kind, &expr)

			require.NoError(t, err)
			require.Len(t, expr.Components, 1)

			component := expr.Components[0]
			assert.Equal(t, tt.expectValue, component.Value)
			assert.Equal(t, tt.weaponName, component.Source)
			assert.True(t, component.Tags.HasTag(tag.FromString(tt.kind)))
		})
	}
}

func TestParseDamageFormula_InvalidFormats(t *testing.T) {
	tests := []struct {
		name    string
		formula string
	}{
		{"empty string", ""},
		{"just letters", "abc"},
		{"float number", "1.5"},
		{"mixed", "1d6+2"},
		{"spaces", "1 d 6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := expression.Expression{Rng: expression.DefaultRoller{}}
			err := parseDamageFormula(tt.formula, "Test", "Damage.Kind.Piercing", &expr)
			assert.Error(t, err)
		})
	}
}

func TestParseDamageFormula_MultipleComponents(t *testing.T) {
	expr := expression.Expression{Rng: expression.DefaultRoller{}}

	// Add dice damage
	err := parseDamageFormula("1d8", "Flaming Sword", "Damage.Kind.Slashing", &expr)
	require.NoError(t, err)

	// Add constant damage
	err = parseDamageFormula("2", "Flaming Sword", "Damage.Kind.Fire", &expr)
	require.NoError(t, err)

	// Should have 2 components
	assert.Len(t, expr.Components, 2)

	// Check first component (dice)
	dice := expr.Components[0]
	assert.Equal(t, 1, dice.Times)
	assert.Equal(t, 8, dice.Sides)
	assert.True(t, dice.Tags.HasTag(tag.FromString("Damage.Kind.Slashing")))

	// Check second component (constant)
	constant := expr.Components[1]
	assert.Equal(t, 2, constant.Value)
	assert.True(t, constant.Tags.HasTag(tag.FromString("Damage.Kind.Fire")))
}

func TestLoadWeapons_Basic(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Create test weapon files
	weaponsDir := tempDir + "/ruleset/weapons"
	err := os.MkdirAll(weaponsDir, 0755)
	require.NoError(t, err)

	// Create test dagger
	daggerYAML := `name: "Test Dagger"
damage:
  - formula: "1d4"
    kind: "Damage.Kind.Piercing"
weapon_tags:
  - "Melee"
  - "Item.Weapon.Simple"
reach: 1`

	err = os.WriteFile(weaponsDir+"/dagger.yml", []byte(daggerYAML), 0644)
	require.NoError(t, err)

	// Create test sword with multiple damage
	swordYAML := `name: "Test Flaming Sword"
damage:
  - formula: "1d8"
    kind: "Damage.Kind.Slashing"
  - formula: "1d6"
    kind: "Damage.Kind.Fire"
weapon_tags:
  - "Melee"
  - "Item.Weapon.Martial"
reach: 1`

	err = os.WriteFile(weaponsDir+"/flamingsword.yml", []byte(swordYAML), 0644)
	require.NoError(t, err)

	// Load weapons
	factories, err := LoadWeapons(tempDir)
	require.NoError(t, err)

	// Check that we got 2 weapons
	assert.Len(t, factories, 2)
	assert.Contains(t, factories, "dagger")
	assert.Contains(t, factories, "flamingsword")

	// Test dagger creation
	dagger := factories["dagger"]()
	assert.Equal(t, "Test Dagger", dagger.Name())
	assert.True(t, dagger.Tags().HasTag(tag.FromString("Melee")))
	assert.True(t, dagger.Tags().HasTag(tag.FromString("Item.Weapon.Simple")))

	// Test dagger damage (1d4 piercing)
	damage := dagger.Damage()
	assert.Len(t, damage.Components, 1)
	assert.Equal(t, 1, damage.Components[0].Times)
	assert.Equal(t, 4, damage.Components[0].Sides)

	// Test flaming sword creation
	sword := factories["flamingsword"]()
	assert.Equal(t, "Test Flaming Sword", sword.Name())
	assert.True(t, sword.Tags().HasTag(tag.FromString("Melee")))
	assert.True(t, sword.Tags().HasTag(tag.FromString("Item.Weapon.Martial")))

	// Test sword damage (1d8 slashing + 1d6 fire)
	swordDamage := sword.Damage()
	assert.Len(t, swordDamage.Components, 2)

	// First component should be 1d8 slashing
	assert.Equal(t, 1, swordDamage.Components[0].Times)
	assert.Equal(t, 8, swordDamage.Components[0].Sides)
	assert.True(t, swordDamage.Components[0].Tags.HasTag(tag.FromString("Damage.Kind.Slashing")))

	// Second component should be 1d6 fire
	assert.Equal(t, 1, swordDamage.Components[1].Times)
	assert.Equal(t, 6, swordDamage.Components[1].Sides)
	assert.True(t, swordDamage.Components[1].Tags.HasTag(tag.FromString("Damage.Kind.Fire")))
}

func TestLoadWeapons_ConstantDamage(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Create test weapon files
	weaponsDir := tempDir + "/ruleset/weapons"
	err := os.MkdirAll(weaponsDir, 0755)
	require.NoError(t, err)

	// Create test weapon with constant damage
	weaponYAML := `name: "Magic Weapon"
damage:
  - formula: "5"
    kind: "Damage.Kind.Force"
weapon_tags:
  - "Melee"
  - "Item.Weapon.Simple"
reach: 1`

	err = os.WriteFile(weaponsDir+"/magic.yml", []byte(weaponYAML), 0644)
	require.NoError(t, err)

	// Load weapons
	factories, err := LoadWeapons(tempDir)
	require.NoError(t, err)

	// Test weapon creation
	weapon := factories["magic"]()
	assert.Equal(t, "Magic Weapon", weapon.Name())

	// Test constant damage
	damage := weapon.Damage()
	assert.Len(t, damage.Components, 1)
	assert.Equal(t, 5, damage.Components[0].Value)
	assert.True(t, damage.Components[0].Tags.HasTag(tag.FromString("Damage.Kind.Force")))
}

func TestLoadWeapons_EmptyDirectory(t *testing.T) {
	// Create a temporary test directory with no weapons
	tempDir := t.TempDir()
	weaponsDir := tempDir + "/ruleset/weapons"
	err := os.MkdirAll(weaponsDir, 0755)
	require.NoError(t, err)

	// Load weapons from empty directory
	factories, err := LoadWeapons(tempDir)
	require.NoError(t, err)

	// Should return empty map
	assert.Empty(t, factories)
}

func TestLoadWeapons_InvalidYAML(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Create test weapon files
	weaponsDir := tempDir + "/ruleset/weapons"
	err := os.MkdirAll(weaponsDir, 0755)
	require.NoError(t, err)

	// Create invalid YAML file
	invalidYAML := `name: "Test Weapon"
damage:
  - formula: "1d4"
    kind: "Damage.Kind.Piercing"
weapon_tags:
  - "Melee"
  - "Item.Weapon.Simple"
reach: invalid_number`

	err = os.WriteFile(weaponsDir+"/invalid.yml", []byte(invalidYAML), 0644)
	require.NoError(t, err)

	// Load weapons should fail
	_, err = LoadWeapons(tempDir)
	assert.Error(t, err)
}
