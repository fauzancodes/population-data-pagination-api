package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Data struct {
	Data [][]interface{} `json:"data"`
}
type DetailKabupatenRequest struct {
	KabID string `json:"id_kab"`
}

func main() {
	app := fiber.New()

	app.Get("/list_kabupaten", func(c *fiber.Ctx) error {
		searchPattern := c.Query("search")
		page := c.QueryInt("page", 1)

		data, err := getData()
		if err != nil {
			log.Fatal(err)
		}

		uniqueNames := make(map[string]map[string]interface{})
		for _, elem := range data.Data {
			id, ok := elem[0].(float64)
			if !ok {
				log.Fatal("Invalid data format")
			}
			name, ok := elem[3].(string)
			if !ok {
				log.Fatal("Invalid data format")
			}
			if _, ok := uniqueNames[name]; !ok {
				uniqueNames[name] = make(map[string]interface{})
				uniqueNames[name]["id_kab"] = id
				uniqueNames[name]["nama"] = name
			}
		}

		var result []map[string]interface{}
		for _, value := range uniqueNames {
			result = append(result, value)
		}
		sort.Slice(result, func(i, j int) bool {
			id1 := result[i]["id_kab"].(float64)
			id2 := result[j]["id_kab"].(float64)
			return id1 < id2
		})
		for i, value := range result {
			value["id_kab"] = i + 1
		}

		var resultPage []map[string]interface{}
		var resultSearch []map[string]interface{}
		limit := 10
		var totalData int
		var totalPage int
		if searchPattern == "" {
			for i, elem := range result {
				if i >= (page-1)*limit && i < page*limit {
					resultPage = append(resultPage, elem)
				}
			}
			totalData = len(result)
			totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
		} else {
			re := regexp.MustCompile(strings.ToLower(searchPattern))
			for _, elem := range result {
				if re.MatchString(strings.ToLower(elem["nama"].(string))) {
					resultSearch = append(resultSearch, elem)
				}
			}
			for i, elem := range resultSearch {
				if i >= (page-1)*limit && i < page*limit {
					resultPage = append(resultPage, elem)
				}
			}
			totalData = len(resultSearch)
			totalPage = int(math.Ceil(float64(totalData) / float64(limit)))
		}
		if page > totalPage {
			return c.JSON(map[string]interface{}{
				"status": false,
				"data":   http.StatusText(http.StatusNotFound),
			})
		} else if page < totalPage {
			return c.JSON(map[string]interface{}{
				"status":       true,
				"data":         resultPage,
				"total":        totalData,
				"search":       searchPattern,
				"limit":        limit,
				"current_page": page,
				"total_page":   totalPage,
				"next":         true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status":       true,
				"data":         resultPage,
				"total":        totalData,
				"search":       searchPattern,
				"limit":        limit,
				"current_page": page,
				"total_page":   totalPage,
				"next":         false,
			})
		}
	})
	app.Post("/detail_kabupaten", func(c *fiber.Ctx) error {
		id := new(DetailKabupatenRequest)
		if err := c.BodyParser(id); err != nil {
			return err
		}
		idInt, _ := strconv.Atoi(id.KabID)

		data, err := getData()
		if err != nil {
			log.Fatal(err)
		}

		uniqueNames := make(map[string]map[string]interface{})
		for _, elem := range data.Data {
			id, ok := elem[0].(float64)
			if !ok {
				log.Fatal("Invalid data format")
			}
			provinsi, ok := elem[2].(string)
			if !ok {
				log.Fatal("Invalid data format")
			}
			name, ok := elem[3].(string)
			if !ok {
				log.Fatal("Invalid data format")
			}
			if _, ok := uniqueNames[name]; !ok {
				uniqueNames[name] = make(map[string]interface{})
				uniqueNames[name]["id_kab"] = id
				uniqueNames[name]["nama_kabupaten"] = name
				uniqueNames[name]["nama_provinsi"] = provinsi
				uniqueNames[name]["total_kecamatan"] = 0
				uniqueNames[name]["total_penduduk"] = 0
			}
		}

		var result []map[string]interface{}
		for _, value := range uniqueNames {
			result = append(result, value)
		}
		sort.Slice(result, func(i, j int) bool {
			id1 := result[i]["id_kab"].(float64)
			id2 := result[j]["id_kab"].(float64)
			return id1 < id2
		})
		for i, value := range result {
			value["id_kab"] = i + 1
		}
		for _, valueData := range result {
			totalKecamatan := 0
			totalPenduduk := 0
			var detailSemuaKecamatan []map[string]interface{}
			for _, item := range data.Data {
				if item[3] == valueData["nama_kabupaten"] {
					totalKecamatan++
					totalPendudukTemp := totalPenduduk + int(item[5].(float64))
					totalPenduduk = totalPendudukTemp
					var detailKecamatan interface{}
					detailKecamatan = map[string]interface{}{
						"id_kec":         item[0],
						"nama_kecamatan": item[4],
						"total":          item[5],
					}
					detailSemuaKecamatan = append(detailSemuaKecamatan, detailKecamatan.(map[string]interface{}))
				}
			}
			for i, value := range detailSemuaKecamatan {
				value["id_kec"] = i + 1
			}
			sort.Slice(detailSemuaKecamatan, func(i, j int) bool {
				total1 := detailSemuaKecamatan[i]["total"].(float64)
				total2 := detailSemuaKecamatan[j]["total"].(float64)
				return total2 < total1
			})
			valueData["total_kecamatan"] = totalKecamatan
			valueData["total_penduduk"] = totalPenduduk
			valueData["list_kecamatan"] = detailSemuaKecamatan
		}

		var resultData map[string]interface{}
		for _, value := range result {
			if idInt == value["id_kab"] {
				resultData = value
			}
		}
		return c.JSON(map[string]interface{}{
			"status": true,
			"data":   resultData,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func getData() (*Data, error) {
	file, err := ioutil.ReadFile("agregatpenduduk2022.json")
	if err != nil {
		return nil, err
	}

	var data Data
	err = json.Unmarshal([]byte(strings.TrimSpace(string(file))), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
