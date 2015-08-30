/*
 * Copyright (c) 2015 Cesanta Software Limited
 * All rights reserved
 */

#include <stdio.h>
#include <string.h>
#include "v7.h"

/*
 * This example demonstrates how to do JS OOP in C.
 */

static v7_val_t MyThing_ctor(struct v7 *v7, v7_val_t this_obj, v7_val_t args) {
  v7_val_t arg0 = v7_array_get(v7, args, 0);
  v7_set(v7, this_obj, "__arg", ~0 /* = strlen */, V7_PROPERTY_DONT_ENUM, arg0);
  return this_obj;
}

static v7_val_t MyThing_myMethod(
    struct v7 *v7, v7_val_t this_obj, v7_val_t args) {
  (void) args;
  return v7_get(v7, this_obj, "__arg", ~0);
}

int main(void) {
  struct v7 *v7 = v7_create();
  v7_val_t ctor_func, proto, eval_result;

  proto = v7_create_object(v7);
  ctor_func = v7_create_constructor(v7, proto, MyThing_ctor, 1);
  v7_set(v7, ctor_func, "MY_CONST", ~0,
         V7_PROPERTY_READ_ONLY | V7_PROPERTY_DONT_DELETE,
         v7_create_number(123));
  v7_set_method(v7, proto, "myMethod", &MyThing_myMethod);
  v7_set(v7, v7_get_global_object(v7), "MyThing", ~0, 0, ctor_func);

  v7_exec(v7, &eval_result, "\
      print('MyThing.MY_CONST = ', MyThing.MY_CONST); \
      var t = new MyThing(456); \
      print('t.MY_CONST = ', t.MY_CONST); \
      print('t.myMethod = ', t.myMethod); \
      print('t.myMethod() = ', t.myMethod());");
  v7_destroy(v7);
  return 0;
}
