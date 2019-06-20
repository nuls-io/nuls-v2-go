package util

import (
	"log"
	"time"
)

var (
	T1 = time.Duration(0)
	T2 = T1
	T3 = T1
	T4 = T1
	T5 = T1
)

func PrintTime() {
	log.Println("t1 : ", T1)
	log.Println("t2 : ", T2)
	log.Println("t3 : ", T3)
	log.Println("t4 : ", T4)
	log.Println("t5 : ", T5)
}