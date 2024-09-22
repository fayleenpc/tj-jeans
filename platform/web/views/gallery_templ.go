// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

var try_script_gallery = templ.NewOnceHandle()

func Gallery(username string) templ.Component {
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
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container-background\" id=\"gallery_page\"><div class=\"slide-background\"><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg1.jpg);\"><div class=\"content-background\"><div class=\"name\">Switzerland</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg2.jpg);\"><div class=\"content-background\"><div class=\"name\">Finland</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg3.jpg);\"><div class=\"content-background\"><div class=\"name\">Iceland</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg4.jpg);\"><div class=\"content-background\"><div class=\"name\">Australia</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg5.jpg);\"><div class=\"content-background\"><div class=\"name\">Netherland</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div><div class=\"item-background\" style=\"background-image: url(/platform/web/static/images/bg6.jpg);\"><div class=\"content-background\"><div class=\"name\">Ireland</div><div class=\"des\">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Ab, eum!</div><button>See More</button></div></div></div><div class=\"button-background\"><button class=\"prev\"><i class=\"fa-solid fa-arrow-left\"></i></button> <button class=\"next\"><i class=\"fa-solid fa-arrow-right\"></i></button></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Var3 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
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
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\r\n        let next = document.querySelector('.next')\r\n        let prev = document.querySelector('.prev')\r\n\r\n\r\n        next.addEventListener('click', function(){\r\n            let items = document.querySelectorAll('.item-background')\r\n            document.querySelector('.slide-background').appendChild(items[0])\r\n        })\r\n\r\n        prev.addEventListener('click', function(){\r\n            let items = document.querySelectorAll('.item-background')\r\n            document.querySelector('.slide-background').prepend(items[items.length - 1]) // here the length of items = 6\r\n        })\r\n    </script>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				return templ_7745c5c3_Err
			})
			templ_7745c5c3_Err = try_script_gallery.Once().Render(templ.WithChildren(ctx, templ_7745c5c3_Var3), templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Page(true, username).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
