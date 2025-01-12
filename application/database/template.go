package database

import "go.mongodb.org/mongo-driver/bson"

func Template(pipeline []bson.D, key string, value interface{}) []bson.D {
	for i, stage := range pipeline {
		pipeline[i] = pipe(stage, key, value)
	}
	return pipeline
}

func pipe(d bson.D, k string, value interface{}) bson.D {
	for i, e := range d {
		switch v := e.Value.(type) {
		case bson.D:
			d[i].Value = pipe(v, k, value)
		case []bson.D:
			d[i].Value = Template(v, k, value)
		default:
			if e.Key == k || e.Key == "$"+k {
				d[i].Value = value
			}
		}
	}
	return d
}
