package cmd

import (
	"log"
	"testing"
)

func TestAddCmd(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    	string
		version    	float32
		want    	[]string
		want2   	bool
		wantErr 	bool
	} {

	}

	for _,v := range  tests {
		log.Println(v)
	}
}
