package models

import "fmt"

type CombatStats struct {
	HP, TempHP, AC, Speed int
}

func (cs *CombatStats) Print() {
	fmt.Printf("HP: %d\n", cs.HP)
	fmt.Printf("Temporary HP: %d\n", cs.TempHP)
	fmt.Printf("AC: %d\n", cs.AC)
	fmt.Printf("Speed: %d\n", cs.Speed)
}
