/* Copyright 2015 the SumatraPDF project authors (see AUTHORS file).
License: Simplified BSD (see COPYING.BSD) */

#include "BaseUtil.h"
#include "WinDynCalls.h"

#define API_DECLARATION(name) \
bool  Has##name = false; \
Sig_##name Dyn##name = nullptr;

KERNEL32_APIS_LIST(API_DECLARATION)
UXTHEME_APIS_LIST(API_DECLARATION)

#undef API_DECLARATION

// Loads a DLL explicitly from the system's library collection
HMODULE SafeLoadLibrary(const WCHAR *dllName) {
    WCHAR dllPath[MAX_PATH];
    UINT res = GetSystemDirectoryW(dllPath, dimof(dllPath));
    if (!res || res >= dimof(dllPath))
        return nullptr;
    BOOL ok = PathAppendW(dllPath, dllName);
    if (!ok)
        return nullptr;
    return LoadLibraryW(dllPath);
}

#define API_LOAD(name) \
    Dyn##name = (Sig_##name)GetProcAddress(h, #name); \
    Has##name = Dyn##name != nullptr;

void InitDynCalls() {
    HMODULE h = SafeLoadLibrary(L"kernel32.dll");
    CrashAlwaysIf(!h);
    KERNEL32_APIS_LIST(API_LOAD);
    h = SafeLoadLibrary(L"uxtheme.dll");
    if (h) {
        UXTHEME_APIS_LIST(API_LOAD);
    }
}

#undef API_LOAD

namespace vss {

    HTHEME OpenThemeData(HWND hwnd, LPCWSTR pszClassList)
    {
        if (HasOpenThemeData) {
            return DynOpenThemeData(hwnd, pszClassList);
        }
        return nullptr;
    }

    HRESULT CloseThemeData(HTHEME hTheme)
    {
        if (HasCloseThemeData) {
            return DynCloseThemeData(hTheme);
        }
        return E_NOTIMPL;
    }

    HRESULT DrawThemeBackground(HTHEME hTheme, HDC hdc, int iPartId, int iStateId, LPCRECT pRect, LPCRECT pClipRect)
    {
        if (HasDrawThemeBackground) {
            return DynDrawThemeBackground(hTheme, hdc, iPartId, iStateId, pRect, pClipRect);
        }
        return E_NOTIMPL;
    }

    BOOL IsThemeActive()
    {
        if (HasIsThemeActive) {
            return DynIsThemeActive();
        }
        return FALSE;
    }

    BOOL IsThemeBackgroundPartiallyTransparent(HTHEME hTheme, int iPartId, int iStateId)
    {
        if (HasIsThemeBackgroundPartiallyTransparent) {
            return DynIsThemeBackgroundPartiallyTransparent(hTheme, iPartId, iStateId);
        }
        return FALSE;
    }

    HRESULT GetThemeColor(HTHEME hTheme, int iPartId, int iStateId, int iPropId, COLORREF *pColor)
    {
        if (HasGetThemeColor) {
            return DynGetThemeColor(hTheme, iPartId, iStateId, iPropId, pColor);
        }
        return E_NOTIMPL;
    }

};
