package docgen

import (
	"strconv"
	"strings"
)

// BaseTemplate is a basic html page with placeholders for: {title}, {css}, {intro}, and {routes}.
func BaseTemplate() string {
	return `
    <html>
      <head>
        <title>{title}</title>
        <style>{css}</style>
        <link rel="icon" type="image/png" href="{favicon.ico}" />
      </head>
      <body>
        <h1>{title}</h1>
        <div>
          {intro}
        </div>
        <div>
          {routes}
        </div>        
      </body>
    </html>
  `
}

// UnorderedList is a list using bullet points instead of numbers.
func UnorderedList(listItems string) string {
	return `
    <ul>
      ` + listItems + `
    </ul>
  `
}

// OrderedList is a numbered list.
func OrderedList(listItems string) string {
	return `
    <ol>
      ` + listItems + `
    </ol>
  `
}

// ListItem is an item for a list.
func ListItem(text string) string {
	return "<li>" + text + "</li>"
}

// Div wraps the string with <div> tags.
func Div(text string) string {
	return "<div>" + text + "</div>"
}

// P wraps the string with <p> tags.
func P(text string) string {
	return "<p>" + text + "</p>"
}

// Head creates a header for a given level eg H1, H2, H3...
func Head(level int, text string) string {
	if len(strings.TrimSpace(text)) == 0 {
		// no text, no header
		return ""
	}
	lvl := strconv.Itoa(level)
	// only H1, H2, H3, H4, H5, H6 are valid.
	// if the level passed in is below 1, use 1
	// if it is above 6, use 6
	if level < 1 {
		lvl = "1"
	} else if level > 6 {
		lvl = "6"
	}
	lvl = "h" + lvl
	header := "<" + lvl + ">" + text + "</" + lvl + ">"

	return header
}

// MilligramMinCSS is a tiny CSS kit, minimized and stringified here.
func MilligramMinCSS() string {
	return `  
    /*!
    * Milligram v1.3.0
    * https://milligram.github.io
    *
    * Copyright (c) 2017 CJ Patoilo
    * Licensed under the MIT license
    */
  
  *,*:after,*:before{box-sizing:inherit}html{box-sizing:border-box;font-size:62.5%}body{color:#606c76;font-family:'Roboto', 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif;font-size:1.6em;font-weight:300;letter-spacing:.01em;line-height:1.6}blockquote{border-left:0.3rem solid #d1d1d1;margin-left:0;margin-right:0;padding:1rem 1.5rem}blockquote *:last-child{margin-bottom:0}.button,button,input[type='button'],input[type='reset'],input[type='submit']{background-color:#9b4dca;border:0.1rem solid #9b4dca;border-radius:.4rem;color:#fff;cursor:pointer;display:inline-block;font-size:1.1rem;font-weight:700;height:3.8rem;letter-spacing:.1rem;line-height:3.8rem;padding:0 3.0rem;text-align:center;text-decoration:none;text-transform:uppercase;white-space:nowrap}.button:focus,.button:hover,button:focus,button:hover,input[type='button']:focus,input[type='button']:hover,input[type='reset']:focus,input[type='reset']:hover,input[type='submit']:focus,input[type='submit']:hover{background-color:#606c76;border-color:#606c76;color:#fff;outline:0}.button[disabled],button[disabled],input[type='button'][disabled],input[type='reset'][disabled],input[type='submit'][disabled]{cursor:default;opacity:.5}.button[disabled]:focus,.button[disabled]:hover,button[disabled]:focus,button[disabled]:hover,input[type='button'][disabled]:focus,input[type='button'][disabled]:hover,input[type='reset'][disabled]:focus,input[type='reset'][disabled]:hover,input[type='submit'][disabled]:focus,input[type='submit'][disabled]:hover{background-color:#9b4dca;border-color:#9b4dca}.button.button-outline,button.button-outline,input[type='button'].button-outline,input[type='reset'].button-outline,input[type='submit'].button-outline{background-color:transparent;color:#9b4dca}.button.button-outline:focus,.button.button-outline:hover,button.button-outline:focus,button.button-outline:hover,input[type='button'].button-outline:focus,input[type='button'].button-outline:hover,input[type='reset'].button-outline:focus,input[type='reset'].button-outline:hover,input[type='submit'].button-outline:focus,input[type='submit'].button-outline:hover{background-color:transparent;border-color:#606c76;color:#606c76}.button.button-outline[disabled]:focus,.button.button-outline[disabled]:hover,button.button-outline[disabled]:focus,button.button-outline[disabled]:hover,input[type='button'].button-outline[disabled]:focus,input[type='button'].button-outline[disabled]:hover,input[type='reset'].button-outline[disabled]:focus,input[type='reset'].button-outline[disabled]:hover,input[type='submit'].button-outline[disabled]:focus,input[type='submit'].button-outline[disabled]:hover{border-color:inherit;color:#9b4dca}.button.button-clear,button.button-clear,input[type='button'].button-clear,input[type='reset'].button-clear,input[type='submit'].button-clear{background-color:transparent;border-color:transparent;color:#9b4dca}.button.button-clear:focus,.button.button-clear:hover,button.button-clear:focus,button.button-clear:hover,input[type='button'].button-clear:focus,input[type='button'].button-clear:hover,input[type='reset'].button-clear:focus,input[type='reset'].button-clear:hover,input[type='submit'].button-clear:focus,input[type='submit'].button-clear:hover{background-color:transparent;border-color:transparent;color:#606c76}.button.button-clear[disabled]:focus,.button.button-clear[disabled]:hover,button.button-clear[disabled]:focus,button.button-clear[disabled]:hover,input[type='button'].button-clear[disabled]:focus,input[type='button'].button-clear[disabled]:hover,input[type='reset'].button-clear[disabled]:focus,input[type='reset'].button-clear[disabled]:hover,input[type='submit'].button-clear[disabled]:focus,input[type='submit'].button-clear[disabled]:hover{color:#9b4dca}code{background:#f4f5f6;border-radius:.4rem;font-size:86%;margin:0 .2rem;padding:.2rem .5rem;white-space:nowrap}pre{background:#f4f5f6;border-left:0.3rem solid #9b4dca;overflow-y:hidden}pre>code{border-radius:0;display:block;padding:1rem 1.5rem;white-space:pre}hr{border:0;border-top:0.1rem solid #f4f5f6;margin:3.0rem 0}input[type='email'],input[type='number'],input[type='password'],input[type='search'],input[type='tel'],input[type='text'],input[type='url'],textarea,select{-webkit-appearance:none;-moz-appearance:none;appearance:none;background-color:transparent;border:0.1rem solid #d1d1d1;border-radius:.4rem;box-shadow:none;box-sizing:inherit;height:3.8rem;padding:.6rem 1.0rem;width:100%}input[type='email']:focus,input[type='number']:focus,input[type='password']:focus,input[type='search']:focus,input[type='tel']:focus,input[type='text']:focus,input[type='url']:focus,textarea:focus,select:focus{border-color:#9b4dca;outline:0}select{background:url('data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' height='14' viewBox='0 0 29 14' width='29'><path fill='#d1d1d1' d='M9.37727 3.625l5.08154 6.93523L19.54036 3.625'/></svg>') center right no-repeat;padding-right:3.0rem}select:focus{background-image:url('data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' height='14' viewBox='0 0 29 14' width='29'><path fill='#9b4dca' d='M9.37727 3.625l5.08154 6.93523L19.54036 3.625'/></svg>')}textarea{min-height:6.5rem}label,legend{display:block;font-size:1.6rem;font-weight:700;margin-bottom:.5rem}fieldset{border-width:0;padding:0}input[type='checkbox'],input[type='radio']{display:inline}.label-inline{display:inline-block;font-weight:normal;margin-left:.5rem}.container{margin:0 auto;max-width:112.0rem;padding:0 2.0rem;position:relative;width:100%}.row{display:flex;flex-direction:column;padding:0;width:100%}.row.row-no-padding{padding:0}.row.row-no-padding>.column{padding:0}.row.row-wrap{flex-wrap:wrap}.row.row-top{align-items:flex-start}.row.row-bottom{align-items:flex-end}.row.row-center{align-items:center}.row.row-stretch{align-items:stretch}.row.row-baseline{align-items:baseline}.row .column{display:block;flex:1 1 auto;margin-left:0;max-width:100%;width:100%}.row .column.column-offset-10{margin-left:10%}.row .column.column-offset-20{margin-left:20%}.row .column.column-offset-25{margin-left:25%}.row .column.column-offset-33,.row .column.column-offset-34{margin-left:33.3333%}.row .column.column-offset-50{margin-left:50%}.row .column.column-offset-66,.row .column.column-offset-67{margin-left:66.6666%}.row .column.column-offset-75{margin-left:75%}.row .column.column-offset-80{margin-left:80%}.row .column.column-offset-90{margin-left:90%}.row .column.column-10{flex:0 0 10%;max-width:10%}.row .column.column-20{flex:0 0 20%;max-width:20%}.row .column.column-25{flex:0 0 25%;max-width:25%}.row .column.column-33,.row .column.column-34{flex:0 0 33.3333%;max-width:33.3333%}.row .column.column-40{flex:0 0 40%;max-width:40%}.row .column.column-50{flex:0 0 50%;max-width:50%}.row .column.column-60{flex:0 0 60%;max-width:60%}.row .column.column-66,.row .column.column-67{flex:0 0 66.6666%;max-width:66.6666%}.row .column.column-75{flex:0 0 75%;max-width:75%}.row .column.column-80{flex:0 0 80%;max-width:80%}.row .column.column-90{flex:0 0 90%;max-width:90%}.row .column .column-top{align-self:flex-start}.row .column .column-bottom{align-self:flex-end}.row .column .column-center{-ms-grid-row-align:center;align-self:center}@media (min-width: 40rem){.row{flex-direction:row;margin-left:-1.0rem;width:calc(100% + 2.0rem)}.row .column{margin-bottom:inherit;padding:0 1.0rem}}a{color:#9b4dca;text-decoration:none}a:focus,a:hover{color:#606c76}dl,ol,ul{list-style:none;margin-top:0;padding-left:0}dl dl,dl ol,dl ul,ol dl,ol ol,ol ul,ul dl,ul ol,ul ul{font-size:90%;margin:1.5rem 0 1.5rem 3.0rem}ol{list-style:decimal inside}ul{list-style:circle inside}.button,button,dd,dt,li{margin-bottom:1.0rem}fieldset,input,select,textarea{margin-bottom:1.5rem}blockquote,dl,figure,form,ol,p,pre,table,ul{margin-bottom:2.5rem}table{border-spacing:0;width:100%}td,th{border-bottom:0.1rem solid #e1e1e1;padding:1.2rem 1.5rem;text-align:left}td:first-child,th:first-child{padding-left:0}td:last-child,th:last-child{padding-right:0}b,strong{font-weight:bold}p{margin-top:0}h1,h2,h3,h4,h5,h6{font-weight:300;letter-spacing:-.1rem;margin-bottom:2.0rem;margin-top:0}h1{font-size:4.6rem;line-height:1.2}h2{font-size:3.6rem;line-height:1.25}h3{font-size:2.8rem;line-height:1.3}h4{font-size:2.2rem;letter-spacing:-.08rem;line-height:1.35}h5{font-size:1.8rem;letter-spacing:-.05rem;line-height:1.5}h6{font-size:1.6rem;letter-spacing:0;line-height:1.4}img{max-width:100%}.clearfix:after{clear:both;content:' ';display:table/*# sourceMappingURL=milligram.min.css.map */}.float-left{float:left}.float-right{float:right}
  `
}

// BassCSS is a zero config drop in css kit.
func BassCSS() string {
	return `
    *{box-sizing:border-box}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,'Helvetica Neue',Helvetica,sans-serif;line-height:1.5;margin:0;color:#111;background-color:#fff}img{max-width:100%;height:auto}svg{max-height:100%}a{color:#07c}h1,h2,h3,h4,h5,h6{font-weight:600;line-height:1.25;margin-top:1em;margin-bottom:.5em}h1{font-size:2rem}h2{font-size:1.5rem}h3{font-size:1.25rem}h4{font-size:1rem}h5{font-size:.875rem}h6{font-size:.75rem}code,pre,samp{font-size:87.5%}blockquote,dl,ol,p,pre,ul{margin-top:1em;margin-bottom:1em}code,pre,samp{font-family:'Roboto Mono','Source Code Pro',Menlo,Consolas,'Liberation Mono',monospace}code,samp{padding:.125em}pre{overflow:scroll}blockquote{font-size:1.25rem;font-style:italic;margin-left:0}hr{margin-top:1.5em;margin-bottom:1.5em;border:0;border-bottom-width:1px;border-bottom-style:solid;border-bottom-color:#ccc}
  `
}

// FaviconIcoData contains image data converted with https://www.motobit.com/util/base64/image-to-base64/
func FaviconIcoData() string {
	return "data:image/x-icon;base64,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Qf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/1tbW/////////////////////////////////////////////////////////////////7+/v/8pKSn/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/DAwM/5GRkf/9/f3///////////////////////////////////////////////////////////////////////b29v9aWlr/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/ycnJ//Z2dn////////////////////////////8/Pz/2NjY/7S0tP+hoaH/qKio/9vb2/////////////////////////////39/f9eXl7/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8WFhb/5+fn/////////////////////////////////1dXV/8AAAD/AAAA/wAAAP8AAAD/AAAA/0dHR//e3t7///////////////////////j4+P80NDT/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/0ZGRv+/v7//i4uL/0xMTP/IyMj/////////////////cnJy/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/xMTE//Nzc3//////////////////////9HR0f8DAwP/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/5ubm/////////////////+YmJj/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/xkZGf/p6en//////////////////////1tbW/8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/bGxs/////////////////8LCwv8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/2FhYf//////////////////////xcXF/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP80NDT/////////////////8PDw/wEBAf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AwMD/+Hh4f/////////////////+/v7/Ghoa/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wMDA//v7+//////////////////JiYm/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/iYmJ//////////////////////9aWlr/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/6qqqv////////////////9cXFz/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP9ISEj//////////////////////4iIiP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/WVlZ/////////////////5WVlf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/x0dHf//////////////////////pqam/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8MDAz/9PT0////////////0NDQ/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/BwcH//////////////////////+0tLT/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP+hoaH////////////8/Pz/DAwM/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8CAgL//////////////////////7Gxsf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/zw8PP////////////////8+Pj7/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/w4ODv//////////////////////oaGh/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/83Nzf///////////21tbf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/LS0t//////////////////////+AgID/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/Wlpa////////////ioqK/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP9lZWX//////////////////////0tLS/8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8CAgL/vLy8//////9wcHD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/8DAwP/////////////////y8vL/CQkJ/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/NjY2/5qamv9xcXH/IiIi/wAAAP8JCQn/Y2Nj/xEREf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP9TU1P//////////////////////5CQkP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8uLi7/+Pj4///////+/v7/zMzM/3x8fP8yMjL/AQEB/wAAAP8AAAD/AAAA/wAAAP8BAQH/Xl5e//X19f/////////////////o6Oj/ExMT/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP9ISEj/7+/v///////////////////////s7Oz/tbW1/42Njf+CgoL/m5ub/9zc3P//////////////////////9fX1/z09Pf8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8dHR3/ra2t//7+/v///////////////////////////////////////////////////////////+Li4v85OTn/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/MjIy/6Ghof/09PT//////////////////////////////////////+fn5/+AgID/DQ0N/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wkJCf9FRUX/d3d3/5OTk/+ZmZn/jIyM/21tbf85OTn/AwMD/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAP8AAAD/AAAA/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
}
