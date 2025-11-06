#!/bin/bash

# Setup test fixtures for MKDIRagons
# This script creates the testdata directory and saves the real API responses

set -e

TESTDATA_DIR="internal/core/testdata"

echo "Creating testdata directory..."
mkdir -p "$TESTDATA_DIR"

echo "Creating wizard.json..."
cat > "$TESTDATA_DIR/wizard.json" << 'EOF'
{
  "index": "wizard",
  "name": "Wizard",
  "hit_die": 6,
  "proficiency_choices": [
    {
      "desc": "Choose two from Arcana, History, Insight, Investigation, Medicine, and Religion",
      "choose": 2,
      "type": "proficiencies",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "reference",
            "item": {
              "index": "skill-arcana",
              "name": "Skill: Arcana",
              "url": "/api/2014/proficiencies/skill-arcana"
            }
          },
          {
            "option_type": "reference",
            "item": {
              "index": "skill-history",
              "name": "Skill: History",
              "url": "/api/2014/proficiencies/skill-history"
            }
          },
          {
            "option_type": "reference",
            "item": {
              "index": "skill-insight",
              "name": "Skill: Insight",
              "url": "/api/2014/proficiencies/skill-insight"
            }
          },
          {
            "option_type": "reference",
            "item": {
              "index": "skill-investigation",
              "name": "Skill: Investigation",
              "url": "/api/2014/proficiencies/skill-investigation"
            }
          },
          {
            "option_type": "reference",
            "item": {
              "index": "skill-medicine",
              "name": "Skill: Medicine",
              "url": "/api/2014/proficiencies/skill-medicine"
            }
          },
          {
            "option_type": "reference",
            "item": {
              "index": "skill-religion",
              "name": "Skill: Religion",
              "url": "/api/2014/proficiencies/skill-religion"
            }
          }
        ]
      }
    }
  ],
  "proficiencies": [
    {
      "index": "daggers",
      "name": "Daggers",
      "url": "/api/2014/proficiencies/daggers"
    },
    {
      "index": "darts",
      "name": "Darts",
      "url": "/api/2014/proficiencies/darts"
    },
    {
      "index": "slings",
      "name": "Slings",
      "url": "/api/2014/proficiencies/slings"
    },
    {
      "index": "quarterstaffs",
      "name": "Quarterstaffs",
      "url": "/api/2014/proficiencies/quarterstaffs"
    },
    {
      "index": "crossbows-light",
      "name": "Crossbows, light",
      "url": "/api/2014/proficiencies/crossbows-light"
    },
    {
      "index": "saving-throw-int",
      "name": "Saving Throw: INT",
      "url": "/api/2014/proficiencies/saving-throw-int"
    },
    {
      "index": "saving-throw-wis",
      "name": "Saving Throw: WIS",
      "url": "/api/2014/proficiencies/saving-throw-wis"
    }
  ],
  "saving_throws": [
    {
      "index": "int",
      "name": "INT",
      "url": "/api/2014/ability-scores/int"
    },
    {
      "index": "wis",
      "name": "WIS",
      "url": "/api/2014/ability-scores/wis"
    }
  ],
  "starting_equipment": [
    {
      "equipment": {
        "index": "spellbook",
        "name": "Spellbook",
        "url": "/api/2014/equipment/spellbook"
      },
      "quantity": 1
    }
  ],
  "starting_equipment_options": [
    {
      "desc": "(a) a quarterstaff or (b) a dagger",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "quarterstaff",
              "name": "Quarterstaff",
              "url": "/api/2014/equipment/quarterstaff"
            }
          },
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "dagger",
              "name": "Dagger",
              "url": "/api/2014/equipment/dagger"
            }
          }
        ]
      }
    },
    {
      "desc": "(a) a component pouch or (b) an arcane focus",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "component-pouch",
              "name": "Component pouch",
              "url": "/api/2014/equipment/component-pouch"
            }
          },
          {
            "option_type": "choice",
            "choice": {
              "desc": "arcane focus",
              "choose": 1,
              "type": "equipment",
              "from": {
                "option_set_type": "equipment_category",
                "equipment_category": {
                  "index": "arcane-foci",
                  "name": "Arcane Foci",
                  "url": "/api/2014/equipment-categories/arcane-foci"
                }
              }
            }
          }
        ]
      }
    },
    {
      "desc": "(a) a scholar's pack or (b) an explorer's pack",
      "choose": 1,
      "type": "equipment",
      "from": {
        "option_set_type": "options_array",
        "options": [
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "scholars-pack",
              "name": "Scholar's Pack",
              "url": "/api/2014/equipment/scholars-pack"
            }
          },
          {
            "option_type": "counted_reference",
            "count": 1,
            "of": {
              "index": "explorers-pack",
              "name": "Explorer's Pack",
              "url": "/api/2014/equipment/explorers-pack"
            }
          }
        ]
      }
    }
  ],
  "class_levels": "/api/2014/classes/wizard/levels",
  "multi_classing": {
    "prerequisites": [
      {
        "ability_score": {
          "index": "int",
          "name": "INT",
          "url": "/api/2014/ability-scores/int"
        },
        "minimum_score": 13
      }
    ],
    "proficiencies": []
  },
  "subclasses": [
    {
      "index": "evocation",
      "name": "Evocation",
      "url": "/api/2014/subclasses/evocation"
    }
  ],
  "spellcasting": {
    "level": 1,
    "spellcasting_ability": {
      "index": "int",
      "name": "INT",
      "url": "/api/2014/ability-scores/int"
    },
    "info": [
      {
        "name": "Cantrips",
        "desc": [
          "At 1st level, you know three cantrips of your choice from the wizard spell list. You learn additional wizard cantrips of your choice at higher levels, as shown in the Cantrips Known column of the Wizard table."
        ]
      },
      {
        "name": "Spellbook",
        "desc": [
          "At 1st level, you have a spellbook containing six 1st- level wizard spells of your choice. Your spellbook is the repository of the wizard spells you know, except your cantrips, which are fixed in your mind."
        ]
      },
      {
        "name": "Preparing and Casting Spells",
        "desc": [
          "The Wizard table shows how many spell slots you have to cast your spells of 1st level and higher. To cast one of these spells, you must expend a slot of the spell's level or higher. You regain all expended spell slots when you finish a long rest.",
          "You prepare the list of wizard spells that are available for you to cast. To do so, choose a number of wizard spells from your spellbook equal to your Intelligence modifier + your wizard level (minimum of one spell). The spells must be of a level for which you have spell slots.",
          "For example, if you're a 3rd-level wizard, you have four 1st-level and two 2nd-level spell slots. With an Intelligence of 16, your list of prepared spells can include six spells of 1st or 2nd level, in any combination, chosen from your spellbook. If you prepare the 1st-level spell magic missile, you can cast it using a 1st-level or a 2nd-level slot. Casting the spell doesn't remove it from your list of prepared spells.",
          "You can change your list of prepared spells when you finish a long rest. Preparing a new list of wizard spells requires time spent studying your spellbook and memorizing the incantations and gestures you must make to cast the spell: at least 1 minute per spell level for each spell on your list."
        ]
      },
      {
        "name": "Spellcasting Ability",
        "desc": [
          "Intelligence is your spellcasting ability for your wizard spells, since you learn your spells through dedicated study and memorization. You use your Intelligence whenever a spell refers to your spellcasting ability. In addition, you use your Intelligence modifier when setting the saving throw DC for a wizard spell you cast and when making an attack roll with one.",
          "Spell save DC = 8 + your proficiency bonus + your Intelligence modifier.",
          "Spell attack modifier = your proficiency bonus + your Intelligence modifier."
        ]
      },
      {
        "name": "Ritual Casting",
        "desc": [
          "You can cast a wizard spell as a ritual if that spell has the ritual tag and you have the spell in your spellbook. You don't need to have the spell prepared."
        ]
      },
      {
        "name": "Spellcasting Focus",
        "desc": [
          "You can use an arcane focus as a spellcasting focus for your wizard spells."
        ]
      }
    ]
  },
  "spells": "/api/2014/classes/wizard/spells",
  "url": "/api/2014/classes/wizard"
}
EOF

echo "Creating dwarf.json..."
cat > "$TESTDATA_DIR/dwarf.json" << 'EOF'
{
  "index": "dwarf",
  "name": "Dwarf",
  "speed": 25,
  "ability_bonuses": [
    {
      "ability_score": {
        "index": "con",
        "name": "CON",
        "url": "/api/2014/ability-scores/con"
      },
      "bonus": 2
    }
  ],
  "alignment": "Most dwarves are lawful, believing firmly in the benefits of a well-ordered society. They tend toward good as well, with a strong sense of fair play and a belief that everyone deserves to share in the benefits of a just order.",
  "age": "Dwarves mature at the same rate as humans, but they're considered young until they reach the age of 50. On average, they live about 350 years.",
  "size": "Medium",
  "size_description": "Dwarves stand between 4 and 5 feet tall and average about 150 pounds. Your size is Medium.",
  "languages": [
    {
      "index": "common",
      "name": "Common",
      "url": "/api/2014/languages/common"
    },
    {
      "index": "dwarvish",
      "name": "Dwarvish",
      "url": "/api/2014/languages/dwarvish"
    }
  ],
  "language_desc": "You can speak, read, and write Common and Dwarvish. Dwarvish is full of hard consonants and guttural sounds, and those characteristics spill over into whatever other language a dwarf might speak.",
  "traits": [
    {
      "index": "darkvision",
      "name": "Darkvision",
      "url": "/api/2014/traits/darkvision"
    },
    {
      "index": "dwarven-resilience",
      "name": "Dwarven Resilience",
      "url": "/api/2014/traits/dwarven-resilience"
    },
    {
      "index": "stonecunning",
      "name": "Stonecunning",
      "url": "/api/2014/traits/stonecunning"
    },
    {
      "index": "dwarven-combat-training",
      "name": "Dwarven Combat Training",
      "url": "/api/2014/traits/dwarven-combat-training"
    },
    {
      "index": "tool-proficiency",
      "name": "Tool Proficiency",
      "url": "/api/2014/traits/tool-proficiency"
    }
  ],
  "subraces": [
    {
      "index": "hill-dwarf",
      "name": "Hill Dwarf",
      "url": "/api/2014/subraces/hill-dwarf"
    }
  ],
  "url": "/api/2014/races/dwarf"
}
EOF

echo "Creating fireball.json..."
cat > "$TESTDATA_DIR/fireball.json" << 'EOF'
{
  "index": "fireball",
  "name": "Fireball",
  "desc": [
    "A bright streak flashes from your pointing finger to a point you choose within range and then blossoms with a low roar into an explosion of flame. Each creature in a 20-foot-radius sphere centered on that point must make a dexterity saving throw. A target takes 8d6 fire damage on a failed save, or half as much damage on a successful one.",
    "The fire spreads around corners. It ignites flammable objects in the area that aren't being worn or carried."
  ],
  "higher_level": [
    "When you cast this spell using a spell slot of 4th level or higher, the damage increases by 1d6 for each slot level above 3rd."
  ],
  "range": "150 feet",
  "components": [
    "V",
    "S",
    "M"
  ],
  "material": "A tiny ball of bat guano and sulfur.",
  "ritual": false,
  "duration": "Instantaneous",
  "concentration": false,
  "casting_time": "1 action",
  "level": 3,
  "damage": {
    "damage_type": {
      "index": "fire",
      "name": "Fire",
      "url": "/api/2014/damage-types/fire"
    },
    "damage_at_slot_level": {
      "3": "8d6",
      "4": "9d6",
      "5": "10d6",
      "6": "11d6",
      "7": "12d6",
      "8": "13d6",
      "9": "14d6"
    }
  },
  "dc": {
    "dc_type": {
      "index": "dex",
      "name": "DEX",
      "url": "/api/2014/ability-scores/dex"
    },
    "dc_success": "half"
  },
  "area_of_effect": {
    "type": "sphere",
    "size": 20
  },
  "school": {
    "index": "evocation",
    "name": "Evocation",
    "url": "/api/2014/magic-schools/evocation"
  },
  "classes": [
    {
      "index": "sorcerer",
      "name": "Sorcerer",
      "url": "/api/2014/classes/sorcerer"
    },
    {
      "index": "wizard",
      "name": "Wizard",
      "url": "/api/2014/classes/wizard"
    }
  ],
  "subclasses": [
    {
      "index": "lore",
      "name": "Lore",
      "url": "/api/2014/subclasses/lore"
    },
    {
      "index": "fiend",
      "name": "Fiend",
      "url": "/api/2014/subclasses/fiend"
    }
  ],
  "url": "/api/2014/spells/fireball"
}
EOF

echo "Creating padded-armor.json..."
cat > "$TESTDATA_DIR/padded-armor.json" << 'EOF'
{
  "desc": [],
  "special": [],
  "index": "padded-armor",
  "name": "Padded Armor",
  "equipment_category": {
    "index": "armor",
    "name": "Armor",
    "url": "/api/2014/equipment-categories/armor"
  },
  "armor_category": "Light",
  "armor_class": {
    "base": 11,
    "dex_bonus": true
  },
  "str_minimum": 0,
  "stealth_disadvantage": true,
  "weight": 8,
  "cost": {
    "quantity": 5,
    "unit": "gp"
  },
  "url": "/api/2014/equipment/padded-armor",
  "contents": [],
  "properties": []
}
EOF

echo ""
echo "✅ Fixtures created successfully!"
echo ""
echo "Directory structure:"
tree -L 2 "$TESTDATA_DIR" 2>/dev/null || ls -la "$TESTDATA_DIR"
echo ""
echo "You can now run tests with:"
echo "  go test -v ./internal/core -run RealJSON"