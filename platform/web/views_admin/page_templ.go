// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package views_admin

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Page(nav bool, username string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link rel=\"icon\" href=\"/platform/web/static/images/TJJeans.ico\"><title>Responsive Admin Dashboard | Korsat X Parmaga</title><!-- ======= Styles ====== --><link rel=\"stylesheet\" href=\"/platform/web/static_admin/style.css\"><script src=\"https://unpkg.com/htmx.org@2.0.2\"></script><script src=\"https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js\" integrity=\"sha384-MrcW6ZMFYlzcLA8Nl+NtUVF0sA7MsXsP1UyJoMp4YLEuNSfAP+JcXn/tWtIaxVXM\" crossorigin=\"anonymous\"></script></head><body><!-- search feature for every db --><!-- <div class=\"search\">\r\n            <label>\r\n                <input type=\"text\" placeholder=\"Search here\">\r\n                <ion-icon name=\"search-outline\"></ion-icon>\r\n            </label>\r\n        </div> --><div class=\"container\"><!-- =============== Navigation ================ -->")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if nav {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"navigation\"><ul><li><a href=\"/admin\"><span class=\"icon\"><ion-icon name=\"logo-apple\"></ion-icon></span> <span class=\"title\">TJ Jeans</span></a></li><li><a href=\"/admin\"><span class=\"icon\"><ion-icon name=\"home-outline\"></ion-icon></span> <span class=\"title\">Dashboard</span></a></li><li><a href=\"/admin/products\"><span class=\"icon\"><ion-icon name=\"people-outline\"></ion-icon></span> <span class=\"title\">Products</span></a></li><li><a href=\"/admin/customers\"><span class=\"icon\"><ion-icon name=\"chatbubble-outline\"></ion-icon></span> <span class=\"title\">Users</span></a></li><li><a href=\"/admin/orders\"><span class=\"icon\"><ion-icon name=\"help-outline\"></ion-icon></span> <span class=\"title\">Orders</span></a></li><li><a href=\"/admin/order_items\"><span class=\"icon\"><ion-icon name=\"settings-outline\"></ion-icon></span> <span class=\"title\">Order Items</span></a></li><li><a hx-post=\"/service/logout\" hx-swap=\"none\"><span class=\"icon\"><ion-icon name=\"log-out-outline\"></ion-icon></span> <span class=\"title\">Sign Out</span></a></li></ul></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"main\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><!-- =========== Scripts =========  --><script src=\"/platform/web/static_admin/app.js\"></script><!-- ====== ionicons ======= --><script type=\"module\" src=\"https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.esm.js\"></script><script nomodule src=\"https://unpkg.com/ionicons@5.5.2/dist/ionicons/ionicons.js\"></script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
