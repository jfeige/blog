package models

import "reflect"

/**
检测 指定的value值是否存在Array，Slice，或者Map的value中
 */

func InArray(obj interface{},target interface{})bool{

	target_tp := reflect.TypeOf(target)
	target_vl := reflect.ValueOf(target)

	switch target_tp.Kind() {
	case reflect.Array,reflect.Slice:
		for i := 0; i < target_vl.Len();i++{
			if obj == target_vl.Index(i).Interface(){
				return true
			}
		}
	case reflect.Map:
		for _,v := range target_vl.MapKeys(){
			if obj == target_vl.MapIndex(v).Interface(){
				return true
			}
		}
	}

	return false
}