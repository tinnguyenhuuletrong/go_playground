package gob_play

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type pokemon struct {
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Skils []skill `json:"skill"`
}

func newPoke(name string, image string, skills []skill) *pokemon {
	return &pokemon{
		Name:  name,
		Image: image,
		Skils: skills,
	}
}

type skill struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func newSkill(name string, image string) *skill {
	return &skill{
		Name:  name,
		Image: image,
	}
}

func Play_Gob() {
	pikachu := newPoke("Pikachu", "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png", []skill{
		*newSkill("せいでんき", ""),
		*newSkill("Static", ""),
	},
	)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	fmt.Println("data:", pikachu)

	enc.Encode(pikachu)
	fmt.Println("gob encoded(hex):", "0x"+hex.EncodeToString(buffer.Bytes()))
	fmt.Println("gob encoded(base64):", base64.RawStdEncoding.EncodeToString(buffer.Bytes()))
	os.WriteFile("./tmp/pika.gob", buffer.Bytes(), 0644)
	fmt.Println("saved to:", "tmp/pika.gob")

	// Decode
	fmt.Println("loaded from:", "tmp/pika.gob")
	fileBytes, err := os.ReadFile("./tmp/pika.gob")
	if err != nil {
		log.Panic(err)
	}
	reader := bytes.NewReader(fileBytes)
	dec := gob.NewDecoder(reader)

	var restoredPika pokemon
	dec.Decode(&restoredPika)

	fmt.Println("restored:", restoredPika)

	fmt.Print("json: ")
	json.NewEncoder(os.Stdout).Encode(pikachu)
}
