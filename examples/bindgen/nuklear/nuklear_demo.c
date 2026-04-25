#include "nuklear.h"
#include <stdio.h>
#include <stdlib.h>
#ifdef _WIN32
#include <windows.h>
#endif

nk_context* nk_init(void) {
    nk_context* ctx = (nk_context*)calloc(1, sizeof(nk_context));
    if (ctx) {
        ctx->frame = 0;
        ctx->button_presses = 0;
        ctx->display_value = 0;
        ctx->pending_value = 0;
        ctx->pending_op = 0;
    }
    printf("[nk] init\n");
    return ctx;
}

int nk_begin(nk_context* ctx, const char* title, int x, int y, int w, int h) {
    if (!ctx) return 0;
    ctx->frame++;
    printf("[nk] begin frame=%d title=%s rect=(%d,%d,%d,%d)\n", ctx->frame, title ? title : "(null)", x, y, w, h);
    return 1;
}

int nk_button_label(nk_context* ctx, const char* label) {
    if (!ctx) return 0;
    ctx->button_presses++;
    printf("[nk] button pressed: %s (count=%d)\n", label ? label : "(null)", ctx->button_presses);
    return 1;
}

void nk_end(nk_context* ctx) {
    if (!ctx) return;
    printf("[nk] end frame=%d\n", ctx->frame);
}

void nk_shutdown(nk_context* ctx) {
    if (!ctx) return;
    printf("[nk] shutdown total_presses=%d\n", ctx->button_presses);
    free(ctx);
}

int nk_button_presses(nk_context* ctx) {
    if (!ctx) return 0;
    return ctx->button_presses;
}

int nk_calc_press_number(nk_context* ctx, int n) {
    if (!ctx) return 0;
    if (n < 0 || n > 9) return 0;
    ctx->display_value = ctx->display_value * 10 + n;
    printf("[nk-calc] number %d -> display=%d\n", n, ctx->display_value);
    return ctx->display_value;
}

int nk_calc_press_op(nk_context* ctx, const char* op) {
    if (!ctx || !op || !op[0]) return 0;
    ctx->pending_value = ctx->display_value;
    ctx->display_value = 0;
    ctx->pending_op = op[0];
    printf("[nk-calc] op %c (pending=%d)\n", ctx->pending_op, ctx->pending_value);
    return 1;
}

int nk_calc_press_equals(nk_context* ctx) {
    if (!ctx) return 0;
    int rhs = ctx->display_value;
    int lhs = ctx->pending_value;
    int out = rhs;
    switch (ctx->pending_op) {
        case '+': out = lhs + rhs; break;
        case '-': out = lhs - rhs; break;
        case '*': out = lhs * rhs; break;
        case '/': out = (rhs == 0) ? 0 : (lhs / rhs); break;
        default: break;
    }
    ctx->display_value = out;
    ctx->pending_value = 0;
    ctx->pending_op = 0;
    printf("[nk-calc] equals -> %d\n", out);
    return out;
}

int nk_calc_display(nk_context* ctx) {
    if (!ctx) return 0;
    return ctx->display_value;
}

void nk_wait_for_enter(void) {
    int ch = 0;
    printf("Press Enter to close...\n");
    fflush(stdout);
    while ((ch = getchar()) != '\n' && ch != EOF) {
    }
}

#ifdef _WIN32
static HWND g_editA;
static HWND g_editB;
static HWND g_result;

static LRESULT CALLBACK NkCalcWndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch (msg) {
        case WM_CREATE: {
            CreateWindowA("STATIC", "A:", WS_CHILD | WS_VISIBLE, 20, 20, 20, 20, hwnd, NULL, NULL, NULL);
            g_editA = CreateWindowA("EDIT", "12", WS_CHILD | WS_VISIBLE | WS_BORDER, 50, 20, 120, 24, hwnd, (HMENU)101, NULL, NULL);
            CreateWindowA("STATIC", "B:", WS_CHILD | WS_VISIBLE, 20, 55, 20, 20, hwnd, NULL, NULL, NULL);
            g_editB = CreateWindowA("EDIT", "7", WS_CHILD | WS_VISIBLE | WS_BORDER, 50, 55, 120, 24, hwnd, (HMENU)102, NULL, NULL);
            CreateWindowA("BUTTON", "Add", WS_CHILD | WS_VISIBLE, 190, 20, 100, 28, hwnd, (HMENU)201, NULL, NULL);
            CreateWindowA("BUTTON", "Subtract", WS_CHILD | WS_VISIBLE, 190, 55, 100, 28, hwnd, (HMENU)202, NULL, NULL);
            CreateWindowA("BUTTON", "Multiply", WS_CHILD | WS_VISIBLE, 300, 20, 100, 28, hwnd, (HMENU)203, NULL, NULL);
            CreateWindowA("BUTTON", "Divide", WS_CHILD | WS_VISIBLE, 300, 55, 100, 28, hwnd, (HMENU)204, NULL, NULL);
            g_result = CreateWindowA("STATIC", "Result: 0", WS_CHILD | WS_VISIBLE, 20, 100, 380, 24, hwnd, NULL, NULL, NULL);
            return 0;
        }
        case WM_COMMAND: {
            int id = LOWORD(wParam);
            if (id >= 201 && id <= 204) {
                char aBuf[64] = {0};
                char bBuf[64] = {0};
                GetWindowTextA(g_editA, aBuf, sizeof(aBuf));
                GetWindowTextA(g_editB, bBuf, sizeof(bBuf));
                int a = atoi(aBuf);
                int b = atoi(bBuf);
                int out = 0;
                const char* op = "+";
                switch (id) {
                    case 201: out = a + b; op = "+"; break;
                    case 202: out = a - b; op = "-"; break;
                    case 203: out = a * b; op = "*"; break;
                    case 204: out = (b == 0) ? 0 : (a / b); op = "/"; break;
                }
                char line[128];
                sprintf(line, "Result: %d %s %d = %d", a, op, b, out);
                SetWindowTextA(g_result, line);
                return 0;
            }
            break;
        }
        case WM_DESTROY:
            PostQuitMessage(0);
            return 0;
    }
    return DefWindowProcA(hwnd, msg, wParam, lParam);
}

void nk_run_calculator_gui(void) {
    HINSTANCE h = GetModuleHandleA(NULL);
    WNDCLASSA wc = {0};
    wc.lpfnWndProc = NkCalcWndProc;
    wc.hInstance = h;
    wc.lpszClassName = "CandyNuklearCalcWnd";
    wc.hCursor = LoadCursor(NULL, IDC_ARROW);
    RegisterClassA(&wc);

    HWND hwnd = CreateWindowA(
        "CandyNuklearCalcWnd",
        "Candy Wrapped C GUI Calculator",
        WS_OVERLAPPEDWINDOW | WS_VISIBLE,
        CW_USEDEFAULT, CW_USEDEFAULT, 460, 210,
        NULL, NULL, h, NULL
    );
    if (!hwnd) {
        printf("[nk] failed to create calculator window\n");
        return;
    }
    MSG msg;
    while (GetMessageA(&msg, NULL, 0, 0) > 0) {
        TranslateMessage(&msg);
        DispatchMessageA(&msg);
    }
}
#else
void nk_run_calculator_gui(void) {
    printf("[nk] calculator GUI is currently implemented for Windows in this demo.\n");
}
#endif
