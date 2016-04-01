package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/boltdb/bolt"
)

// type Protein struct {
// 	ID string
// 	//original
// 	Seq    string
// 	Dssp   string
// 	Stride string
// 	Kaksi  string
// 	Pross  string
// 	//processed
// 	Dssp3   string
// 	Stride3 string
// 	Kaksi3  string
// 	Pross3  string
// 	// consensus 2
// 	DsspStride3  string
// 	DsspKaksi3   string
// 	DsspPross3   string
// 	StrideKaksi3 string
// 	StridePross3 string
// 	KaksiPross3  string
// 	// consensus 3
// 	DsspStrideKaksi3  string
// 	DsspStridePross3  string
// 	DsspKaksiPross3   string
// 	StrideKaksiPross3 string
// 	// consensus 4
// 	All3 string
// }

// func readJson(dir, fn string) ([]byte, error) {
// 	// path := "/home/jgcarvalho/sync/data/multissdb/pdb_2A_chain_min80/103lA"
// 	var data Protein
// 	path := dir + fn
// 	file, err := ioutil.ReadFile(path)
// 	lines, _ := ioutil.ReadFile(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = json.Unmarshal(file, &data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(data)
// 	// if len(lines) > 0 {
// 	// 	fields := strings.Fields(string(lines))
// 	// 	protein := &Protein{
// 	// 		ID:       fn,
// 	// 		Seq:      fields[1],
// 	// 		SSdssp:   fields[3],
// 	// 		SSstride: fields[5],
// 	// 		SSpross:  fields[7],
// 	// 		SSkaksi:  fields[9]}
// 	// 	return json.Marshal(protein)
// 	// } else {
// 	// 	return []byte(""), fmt.Errorf("Empty file: %v", fn)
// 	// }
//
// }

func createDB(dbname, dir string) {
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	all, _ := ioutil.ReadDir(dir)
	for i, v := range all {
		fn := v.Name()
		fmt.Println(i, fn)
		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte("proteins"))
			if err != nil {
				return err
			}
			// encoded, err := json.Marshal(protein)
			// encoded, err := file2json(dir, fn)
			// if err != nil {
			// 	return err
			// }
			json, err := ioutil.ReadFile(dir + fn)
			if err != nil {
				log.Fatal(err)
			}
			return b.Put([]byte(fn), json)
		})
	}
}

func viewDB(dbname string) {
	db, err := bolt.Open(dbname, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("proteins"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

}

func main() {
	cr := flag.Bool("create", false, "Create database")
	view := flag.Bool("view", false, "View database")
	dbname := flag.String("dbname", "", "DataBase filename")
	dirname := flag.String("dir", "", "Data path")
	flag.Parse()

	if *cr {
		if *dbname != "" && *dirname != "" {
			createDB(*dbname, *dirname)
		}
	}

	if *view {
		if *dbname != "" {
			viewDB(*dbname)
		}
	}
	// dir := "/home/jgcarvalho/sync/data/multissdb/chameleonic/"
	// dbname := "chameleonic.db"

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.

}
