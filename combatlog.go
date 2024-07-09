package main

import "fmt"

var combatLog []string

func addDamageDone(npcName NpcName, damage int, crit, dodge bool) {
	var text string
	if dodge {
		text = fmt.Sprintf("%s dodged your attack", npcName)
	} else if crit {
		text = fmt.Sprintf("%d damage done to %s (critical)", damage, npcName)
	} else {
		text = fmt.Sprintf("%d damage done to %s", damage, npcName)
	}

	if len(combatLog)+1 > 5 {
		combatLog = append(combatLog[:1], combatLog[2:]...)
	}
	combatLog = append(combatLog, text)
}

func addDamageTaken(npcName NpcName, damage int, crit, dodge bool) {
	var text string
	if dodge {
		text = fmt.Sprintf("You dodged %s's attack", npcName)
	} else if crit {
		text = fmt.Sprintf("You took %d damage from %s (critical)", damage, npcName)
	} else {
		text = fmt.Sprintf("You took %d damage from %s", damage, npcName)
	}

	if len(combatLog)+1 > 5 {
		combatLog = append(combatLog[:1], combatLog[2:]...)
	}
	combatLog = append(combatLog, text)
}

func addPowerupPickup(powerup string, modifier float64) {
	text := fmt.Sprintf("Your %s has increased %.2f%%", powerup, modifier)
	if len(combatLog)+1 > 5 {
		combatLog = append(combatLog[:1], combatLog[2:]...)
	}
	combatLog = append(combatLog, text)
}
