// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/fayleenpc/tj-jeans/platform/web/views/components"

var try_script_page = templ.NewOnceHandle()

func Page(nav bool, username string, role string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>TJ Jeans Store</title><link href=\"https://unpkg.com/boxicons@2.1.4/css/boxicons.min.css\" rel=\"stylesheet\"><link rel=\"stylesheet\" href=\"/platform/web/static/style.css\"><link rel=\"icon\" href=\"/platform/web/static/images/TJJeans.ico\"><script src=\"https://unpkg.com/htmx.org@2.0.2\"></script></head><body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if nav {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"nav\"><div class=\"nav-logo\"><p>TJ Jeans</p></div><div class=\"nav-menu\" id=\"navMenu\"><ul><li><a href=\"/\" class=\"link\" id=\"home\" onclick=\"addActiveClass(document.getElementById(&#39;home&#39;), &#39;active&#39;)\">Beranda</a></li><li><a href=\"/products\" class=\"link\" onclick=\"addActiveClass(document.getElementById(&#39;products&#39;), &#39;active&#39;)\" id=\"products\">Produk</a></li><li><a href=\"/gallery\" class=\"link\" id=\"gallery\" onclick=\"addActiveClass(document.getElementById(&#39;gallery&#39;), &#39;active&#39;)\">Galeri</a></li>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.DisplayMenuLogin(username).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.DisplayMenuAdmin(role).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul></div><div class=\"nav-button\" id=\"nav-button\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.DisplayButtonLogout(username).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"nav-menu-btn\"><i class=\"bx bx-menu\" onclick=\"myMenuFunction()\"></i></div></nav>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"notifications\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\r\n        let notifications = document.querySelector('.notifications');\r\n        function myMenuFunction() {\r\n                var i = document.getElementById(\"navMenu\");\r\n\r\n                if(i.className === \"nav-menu\") {\r\n                    i.className += \" responsive\";\r\n                } else {\r\n                    i.className = \"nav-menu\";\r\n                }\r\n        }\r\n        \r\n        function hasActiveClass(el, className)\r\n        {\r\n            if (el.classList) {\r\n                return el.classList.contains(className);\r\n            }\r\n            return !!el.className.match(new RegExp('(\\\\s|^)' + className + '(\\\\s|$)'));\r\n        }\r\n\r\n        function addActiveClass(el, className)\r\n        {\r\n            if (el.classList)\r\n                el.classList.add(className)\r\n            else if (!hasClass(el, className))\r\n                el.className += \" \" + className;\r\n        }\r\n\r\n        function removeActiveClass(el, className)\r\n        {\r\n            if (el.classList)\r\n                el.classList.remove(className)\r\n            else if (hasClass(el, className))\r\n            {\r\n                var reg = new RegExp('(\\\\s|^)' + className + '(\\\\s|$)');\r\n                el.className = el.className.replace(reg, ' ');\r\n            }\r\n        }\r\n    </script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = try_script_page.Once().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
