package class_test

import (
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// ClassModelTestSuite defines the test suite for Class model
type ClassModelTestSuite struct {
	suite.Suite
	class *class.Class
}

// SetupTest runs before each test
func (suite *ClassModelTestSuite) SetupTest() {
	suite.class = &class.Class{
		Index:  "wizard",
		Name:   "Wizard",
		HitDie: 6,
		URL:    "/api/classes/wizard",
		Proficiencies: []reference.Reference{
			{Index: "daggers", Name: "Daggers", URL: "/api/proficiencies/daggers"},
		},
		SavingThrows: []reference.Reference{
			{Index: "int", Name: "INT", URL: "/api/ability-scores/int"},
			{Index: "wis", Name: "WIS", URL: "/api/ability-scores/wis"},
		},
	}
}

// TestClassBasicFields tests basic field assignment and retrieval
func (suite *ClassModelTestSuite) TestClassBasicFields() {
	assert.Equal(suite.T(), "wizard", suite.class.Index)
	assert.Equal(suite.T(), "Wizard", suite.class.Name)
	assert.Equal(suite.T(), 6, suite.class.HitDie)
	assert.Equal(suite.T(), "/api/classes/wizard", suite.class.URL)
}

// TestClassProficiencies tests proficiencies
func (suite *ClassModelTestSuite) TestClassProficiencies() {
	assert.Len(suite.T(), suite.class.Proficiencies, 1)
	assert.Equal(suite.T(), "daggers", suite.class.Proficiencies[0].Index)
	assert.Equal(suite.T(), "Daggers", suite.class.Proficiencies[0].Name)
}

// TestClassSavingThrows tests saving throws
func (suite *ClassModelTestSuite) TestClassSavingThrows() {
	assert.Len(suite.T(), suite.class.SavingThrows, 2)
	assert.Equal(suite.T(), "int", suite.class.SavingThrows[0].Index)
	assert.Equal(suite.T(), "wis", suite.class.SavingThrows[1].Index)
}

// TestGetEndpoint tests the GetEndpoint method
func (suite *ClassModelTestSuite) TestGetEndpoint() {
	endpoint := suite.class.GetEndpoint()
	assert.Equal(suite.T(), "classes/", endpoint)
}

// TestPrint tests the Print method doesn't panic
func (suite *ClassModelTestSuite) TestPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.class.Print()
	})
}

// TestPrintFeatures tests the PrintFeatures method doesn't panic
func (suite *ClassModelTestSuite) TestPrintFeatures() {
	assert.NotPanics(suite.T(), func() {
		suite.class.PrintFeatures()
	})
}

// Run the test suite
func TestClassModelTestSuite(t *testing.T) {
	suite.Run(t, new(ClassModelTestSuite))
}

// TestProficiencyChoice tests ProficiencyChoice struct
func TestProficiencyChoice(t *testing.T) {
	pc := class.ProficiencyChoice{
		Desc:   "Choose 2 skills",
		Choose: 2,
		Type:   "proficiencies",
		From: class.OptionGroup{
			OptionSetType: "options_array",
			Options: []class.Option{
				{
					OptionType: "reference",
					Item: &reference.Reference{
						Index: "skill-arcana",
						Name:  "Skill: Arcana",
					},
				},
			},
		},
	}

	assert.Equal(t, "Choose 2 skills", pc.Desc)
	assert.Equal(t, 2, pc.Choose)
	assert.Equal(t, "proficiencies", pc.Type)
	assert.Equal(t, "options_array", pc.From.OptionSetType)
	assert.Len(t, pc.From.Options, 1)
	assert.NotNil(t, pc.From.Options[0].Item)
	assert.Equal(t, "skill-arcana", pc.From.Options[0].Item.Index)
}

// TestOptionWithCount tests Option struct with count
func TestOptionWithCount(t *testing.T) {
	opt := class.Option{
		OptionType: "counted_reference",
		Item: &reference.Reference{
			Index: "dagger",
			Name:  "Dagger",
		},
		Count: 2,
	}

	assert.Equal(t, "counted_reference", opt.OptionType)
	assert.NotNil(t, opt.Item)
	assert.Equal(t, 2, opt.Count)
	assert.Nil(t, opt.Of)
	assert.Nil(t, opt.Choice)
}

// TestOptionWithChoice tests Option struct with nested choice
func TestOptionWithChoice(t *testing.T) {
	opt := class.Option{
		OptionType: "choice",
		Choice: &class.ChoiceGroup{
			Desc:   "Select one",
			Choose: 1,
			Type:   "equipment",
			From: class.EquipmentCategory{
				OptionSetType: "equipment_category",
				EquipmentCategory: reference.Reference{
					Index: "simple-weapons",
					Name:  "Simple Weapons",
				},
			},
		},
	}

	assert.Equal(t, "choice", opt.OptionType)
	assert.NotNil(t, opt.Choice)
	assert.Equal(t, "Select one", opt.Choice.Desc)
	assert.Equal(t, 1, opt.Choice.Choose)
	assert.Equal(t, "simple-weapons", opt.Choice.From.EquipmentCategory.Index)
}

// TestStartingEquipment tests StartingEquipment struct
func TestStartingEquipment(t *testing.T) {
	se := class.StartingEquipment{
		Equipment: reference.Reference{
			Index: "leather-armor",
			Name:  "Leather Armor",
			URL:   "/api/equipment/leather-armor",
		},
		Quantity: 1,
	}

	assert.Equal(t, "leather-armor", se.Equipment.Index)
	assert.Equal(t, "Leather Armor", se.Equipment.Name)
	assert.Equal(t, 1, se.Quantity)
}

// TestMultiClassing tests MultiClassing struct
func TestMultiClassing(t *testing.T) {
	mc := class.MultiClassing{
		Prerequisites: []class.Prerequisite{
			{
				AbilityScore: reference.Reference{
					Index: "int",
					Name:  "INT",
				},
				MinimumScore: 13,
			},
		},
		Proficiencies: []interface{}{},
	}

	assert.Len(t, mc.Prerequisites, 1)
	assert.Equal(t, "int", mc.Prerequisites[0].AbilityScore.Index)
	assert.Equal(t, 13, mc.Prerequisites[0].MinimumScore)
	assert.Empty(t, mc.Proficiencies)
}

// TestSpellcasting tests Spellcasting struct
func TestSpellcasting(t *testing.T) {
	sc := class.Spellcasting{
		Level: 1,
		SpellcastingAbility: reference.Reference{
			Index: "int",
			Name:  "INT",
			URL:   "/api/ability-scores/int",
		},
		Info: []class.SpellcastingInfo{
			{
				Name: "Cantrips",
				Desc: []string{"You know three cantrips of your choice from the wizard spell list."},
			},
		},
	}

	assert.Equal(t, 1, sc.Level)
	assert.Equal(t, "int", sc.SpellcastingAbility.Index)
	assert.Len(t, sc.Info, 1)
	assert.Equal(t, "Cantrips", sc.Info[0].Name)
	assert.Len(t, sc.Info[0].Desc, 1)
}

// TestEmptyClass tests empty class initialization
func TestEmptyClass(t *testing.T) {
	testClass := &class.Class{}

	assert.Empty(t, testClass.Index)
	assert.Empty(t, testClass.Name)
	assert.Equal(t, 0, testClass.HitDie)
	assert.Nil(t, testClass.Proficiencies)
	assert.Nil(t, testClass.SavingThrows)
	assert.Equal(t, "classes/", testClass.GetEndpoint())
}

// TestComplexClass tests a fully populated class
func TestComplexClass(t *testing.T) {
	testClass := &class.Class{
		Index:  "fighter",
		Name:   "Fighter",
		HitDie: 10,
		ProficiencyChoices: []class.ProficiencyChoice{
			{
				Desc:   "Choose 2 skills",
				Choose: 2,
				Type:   "proficiencies",
			},
		},
		StartingEquipment: []class.StartingEquipment{
			{
				Equipment: reference.Reference{Index: "chain-mail"},
				Quantity:  1,
			},
		},
		StartingEquipmentOptions: []class.EquipmentOptionGroup{
			{
				Desc:   "Choose 1",
				Choose: 1,
				Type:   "equipment",
			},
		},
		MultiClassing: class.MultiClassing{
			Prerequisites: []class.Prerequisite{
				{
					AbilityScore: reference.Reference{Index: "str"},
					MinimumScore: 13,
				},
			},
		},
	}

	assert.Equal(t, "fighter", testClass.Index)
	assert.Len(t, testClass.ProficiencyChoices, 1)
	assert.Len(t, testClass.StartingEquipment, 1)
	assert.Len(t, testClass.StartingEquipmentOptions, 1)
	assert.Len(t, testClass.MultiClassing.Prerequisites, 1)
}
