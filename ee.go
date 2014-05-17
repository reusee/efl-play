package main

/*
#define EFL_EO_API_SUPPORT
#define EFL_BETA_API_SUPPORT

#include <Ecore.h>
#include <Eo.h>
#include <Evas.h>
#include <Ecore_Evas.h>
#cgo pkg-config: ecore eo evas ecore-evas

void quit(Ecore_Evas *ee) {
	ecore_main_loop_quit();
}

void set_delete_cb(Ecore_Evas *ee) {
	ecore_evas_callback_delete_request_set(ee, quit);
}

*/
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

func main() {
	C.ecore_evas_init()
	ee := C.ecore_evas_new(nil, 0, 0, 800, 600, nil)
	C.set_delete_cb(ee)
	C.ecore_evas_title_set(ee, C.CString("EcoreEvas例子"))
	C.ecore_evas_show(ee)
	canvas := C.ecore_evas_get(ee)

	bg := Add(C.evas_obj_rectangle_class_get(), canvas, func() {
		C.evas_obj_color_set(255, 255, 255, 255)
		C.evas_obj_size_set(800, 600)
		C.evas_obj_visibility_set(C.EINA_TRUE)
	})
	_ = bg

	box := Add(C.evas_obj_rectangle_class_get(), canvas, func() {
		C.evas_obj_color_set(255, 0, 0, 255)
		C.evas_obj_position_set(30, 30)
		C.evas_obj_size_set(100, 100)
		C.evas_obj_visibility_set(C.EINA_TRUE)
	})
	_ = box

	C.ecore_main_loop_begin()
}

func Do(id *C.Eo, f func()) {
	pc, file, line, _ := runtime.Caller(1)
	cFile := C.CString(file)
	cFunc := C.CString(runtime.FuncForPC(pc).Name())
	C._eo_do_start(id, nil, C.EINA_FALSE, cFile, cFunc, C.int(line))
	f()
	C._eo_do_end(&id)
	C.free(unsafe.Pointer(cFile))
	C.free(unsafe.Pointer(cFunc))
}

func Add(class *C.Eo_Class, parent interface{}, fs ...func()) *C.Eo {
	tmpClass := class
	_, file, line, _ := runtime.Caller(1)
	cFile := C.CString(file)
	tmpObj := C._eo_add_internal_start(cFile, C.int(line), tmpClass,
		(*C.Eo)(unsafe.Pointer(reflect.ValueOf(parent).Pointer())))
	Do(tmpObj, func() {
		C.eo_constructor()
		for _, f := range fs {
			f()
		}
		tmpObj = C._eo_add_internal_end(cFile, C.int(line), tmpObj)
	})
	return tmpObj
}
