%any = type { i8, i8* }
declare i64 @candy_dyn_get_i64(i8*, i8*)
declare i64 @candy_dyn_call_i64(i8*, i8*)
declare i32 @printf(i8*, ...)
declare void @exit(i32)
declare i64 @strlen(i8*)
declare i64 @clock()
declare i8* @malloc(i64)
declare i8* @strcpy(i8*, i8*)
declare i8* @strcat(i8*, i8*)
declare i32 @sprintf(i8*, i8*, ...)

define i8* @candy_str_add(i8* %s1, i8* %s2) {
  %l1 = call i64 @strlen(i8* %s1)
  %l2 = call i64 @strlen(i8* %s2)
  %l = add i64 %l1, %l2
  %ltot = add i64 %l, 1
  %ptr = call i8* @malloc(i64 %ltot)
  %tmp1 = call i8* @strcpy(i8* %ptr, i8* %s1)
  %tmp2 = call i8* @strcat(i8* %ptr, i8* %s2)
  ret i8* %ptr
}

define i8* @candy_int_to_str(i64 %val) {
  %ptr = call i8* @malloc(i64 32)
  %fmt = getelementptr inbounds [5 x i8], [5 x i8]* @.str.fmt_int_only, i64 0, i64 0
  %tmp = call i32 (i8*, i8*, ...) @sprintf(i8* %ptr, i8* %fmt, i64 %val)
  ret i8* %ptr
}
@.str.fmt_int = private unnamed_addr constant [6 x i8] c"%lld\0A\00"
@.str.fmt_int_only = private unnamed_addr constant [5 x i8] c"%lld\00"
@.str.fmt_float = private unnamed_addr constant [4 x i8] c"%f\0A\00"
@.str.fmt_str = private unnamed_addr constant [4 x i8] c"%s\0A\00"
%vector2 = type {double, i64, double}

define i64 @main() {
  %1 = alloca %vector2
  %2 = getelementptr inbounds %vector2, %vector2* %1, i32 0, i32 0
  store double 1.000000, double* %2
  %3 = getelementptr inbounds %vector2, %vector2* %1, i32 0, i32 2
  store double 2.000000, double* %3
  %v1 = alloca %vector2*
  store %vector2* %1, %vector2** %v1
  ; compiling return statement
  ret i64 0
}
