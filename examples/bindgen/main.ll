%any = type { i8, i8* }
declare i64 @candy_dyn_get_i64(i8*, i8*)
declare i64 @candy_dyn_call_i64(i8*, i8*)
declare i32 @printf(i8*, ...)
declare void @exit(i32)
declare i64 @strlen(i8*)
declare i64 @clock()
declare double @llvm.sqrt.f64(double)
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
declare i64 @add(i64, i64)
declare i8* @version()
@.str.s1 = private unnamed_addr constant [14 x i8] c"add(10,32) = \00"

define i64 @main() nounwind willreturn hot inlinehint {
  %1 = call i64 @add(i64 10, i64 32)
  %sum = alloca i64
  store i64 %1, i64* %sum
  %2 = load i64, i64* %sum
  %3 = getelementptr inbounds [14 x i8], [14 x i8]* @.str.s1, i64 0, i64 0
  %4 = load i64, i64* %sum
  %5 = call i8* @candy_int_to_str(i64 %4)
  %6 = call i8* @candy_str_add(i8* %3, i8* %5)
  %7 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.fmt_str, i64 0, i64 0), i8* %6)
  ret i64 0
}
