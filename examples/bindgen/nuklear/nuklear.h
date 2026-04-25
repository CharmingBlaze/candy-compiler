// Minimal Nuklear-style C API fixture for candywrap integration tests.
#ifndef NK_MINI_H
#define NK_MINI_H

typedef struct nk_context {
    int frame;
    int button_presses;
    int display_value;
    int pending_value;
    char pending_op;
} nk_context;

nk_context* nk_init(void);
int nk_begin(nk_context* ctx, const char* title, int x, int y, int w, int h);
int nk_button_label(nk_context* ctx, const char* label);
void nk_end(nk_context* ctx);
void nk_shutdown(nk_context* ctx);
int nk_button_presses(nk_context* ctx);
int nk_calc_press_number(nk_context* ctx, int n);
int nk_calc_press_op(nk_context* ctx, const char* op);
int nk_calc_press_equals(nk_context* ctx);
int nk_calc_display(nk_context* ctx);
void nk_wait_for_enter(void);
void nk_run_calculator_gui(void);

#endif
