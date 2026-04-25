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
declare i64 @nk_init()
declare i64 @nk_begin(i64, i8*, i64, i64, i64, i64)
declare i64 @nk_button_label(i64, i8*)
declare i64 @nk_end(i64)
declare i64 @nk_shutdown(i64)
declare i64 @nk_button_presses(i64)
declare i64 @nk_calc_press_number(i64, i64)
declare i64 @nk_calc_press_op(i64, i8*)
declare i64 @nk_calc_press_equals(i64)
declare i64 @nk_calc_display(i64)
declare i64 @nk_wait_for_enter()
declare i64 @nk_run_calculator_gui()
declare i64 @nk_NkCalcWndProc(i64, i64, i64, i64)
declare i64 @nk_DefWindowProcA(i64, i64, i64, i64)
@.str.s1 = private unnamed_addr constant [19 x i8] c"Candy Nuklear Demo\00"
@.str.s3 = private unnamed_addr constant [17 x i8] c"Hello from Candy\00"
@.str.s5 = private unnamed_addr constant [26 x i8] c"GUI event: button clicked\00"
@.str.s6 = private unnamed_addr constant [15 x i8] c"press count = \00"
@.str.s7 = private unnamed_addr constant [2 x i8] c"+\00"
@.str.s8 = private unnamed_addr constant [21 x i8] c"calculator result = \00"
@.str.s9 = private unnamed_addr constant [11 x i8] c" (display=\00"
@.str.s10 = private unnamed_addr constant [2 x i8] c")\00"
@.str.s11 = private unnamed_addr constant [40 x i8] c"Opening native calculator GUI window...\00"

define i64 @main() nounwind willreturn hot inlinehint {
  %1 = call i64 @nk_init()
  %ctx = alloca i64
  store i64 %1, i64* %ctx
  %2 = load i64, i64* %ctx
  %3 = load i64, i64* %ctx
  %4 = getelementptr inbounds [19 x i8], [19 x i8]* @.str.s1, i64 0, i64 0
  %5 = call i64 @nk_begin(i64 %3, i8* %4, i64 20, i64 20, i64 320, i64 200)
  %opened = alloca i64
  store i64 %5, i64* %opened
  %6 = load i64, i64* %opened
  %7 = load i64, i64* %opened
  %8 = icmp eq i64 %7, 1
  br i1 %8, label %then2, label %merge2
then2:
  %9 = load i64, i64* %ctx
  %10 = getelementptr inbounds [17 x i8], [17 x i8]* @.str.s3, i64 0, i64 0
  %11 = call i64 @nk_button_label(i64 %9, i8* %10)
  %clicked = alloca i64
  store i64 %11, i64* %clicked
  %12 = load i64, i64* %clicked
  %13 = load i64, i64* %clicked
  %14 = icmp eq i64 %13, 1
  br i1 %14, label %then4, label %merge4
then4:
  %15 = getelementptr inbounds [26 x i8], [26 x i8]* @.str.s5, i64 0, i64 0
  %16 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.fmt_str, i64 0, i64 0), i8* %15)
  br label %merge4
merge4:
  br label %merge2
merge2:
  %17 = load i64, i64* %ctx
  %18 = call i64 @nk_end(i64 %17)
  %19 = load i64, i64* %ctx
  %20 = call i64 @nk_button_presses(i64 %19)
  %presses = alloca i64
  store i64 %20, i64* %presses
  %21 = load i64, i64* %presses
  %22 = getelementptr inbounds [15 x i8], [15 x i8]* @.str.s6, i64 0, i64 0
  %23 = load i64, i64* %presses
  %24 = call i8* @candy_int_to_str(i64 %23)
  %25 = call i8* @candy_str_add(i8* %22, i8* %24)
  %26 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.fmt_str, i64 0, i64 0), i8* %25)
  %27 = load i64, i64* %ctx
  %28 = call i64 @nk_calc_press_number(i64 %27, i64 1)
  %29 = load i64, i64* %ctx
  %30 = call i64 @nk_calc_press_number(i64 %29, i64 2)
  %31 = load i64, i64* %ctx
  %32 = getelementptr inbounds [2 x i8], [2 x i8]* @.str.s7, i64 0, i64 0
  %33 = call i64 @nk_calc_press_op(i64 %31, i8* %32)
  %34 = load i64, i64* %ctx
  %35 = call i64 @nk_calc_press_number(i64 %34, i64 7)
  %36 = load i64, i64* %ctx
  %37 = call i64 @nk_calc_press_equals(i64 %36)
  %result = alloca i64
  store i64 %37, i64* %result
  %38 = load i64, i64* %result
  %39 = load i64, i64* %ctx
  %40 = call i64 @nk_calc_display(i64 %39)
  %display = alloca i64
  store i64 %40, i64* %display
  %41 = load i64, i64* %display
  %42 = getelementptr inbounds [21 x i8], [21 x i8]* @.str.s8, i64 0, i64 0
  %43 = load i64, i64* %result
  %44 = call i8* @candy_int_to_str(i64 %43)
  %45 = call i8* @candy_str_add(i8* %42, i8* %44)
  %46 = getelementptr inbounds [11 x i8], [11 x i8]* @.str.s9, i64 0, i64 0
  %47 = call i8* @candy_str_add(i8* %45, i8* %46)
  %48 = load i64, i64* %display
  %49 = call i8* @candy_int_to_str(i64 %48)
  %50 = call i8* @candy_str_add(i8* %47, i8* %49)
  %51 = getelementptr inbounds [2 x i8], [2 x i8]* @.str.s10, i64 0, i64 0
  %52 = call i8* @candy_str_add(i8* %50, i8* %51)
  %53 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.fmt_str, i64 0, i64 0), i8* %52)
  %54 = load i64, i64* %ctx
  %55 = call i64 @nk_shutdown(i64 %54)
  %56 = getelementptr inbounds [40 x i8], [40 x i8]* @.str.s11, i64 0, i64 0
  %57 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.fmt_str, i64 0, i64 0), i8* %56)
  %58 = call i64 @nk_run_calculator_gui()
  ret i64 0
}
