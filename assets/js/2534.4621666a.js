(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[2534],{422:(e,t,n)=>{"use strict";n.d(t,{A:()=>x});n(8101);var s=n(3526),o=n(4323),c=n(6435),a=n(9515);const r={iconEdit:"iconEdit_VA3U"};var i=n(5105);function l(e){let{className:t,...n}=e;return(0,i.jsx)("svg",{fill:"currentColor",height:"20",width:"20",viewBox:"0 0 40 40",className:(0,s.A)(r.iconEdit,t),"aria-hidden":"true",...n,children:(0,i.jsx)("g",{children:(0,i.jsx)("path",{d:"m34.5 11.7l-3 3.1-6.3-6.3 3.1-3q0.5-0.5 1.2-0.5t1.1 0.5l3.9 3.9q0.5 0.4 0.5 1.1t-0.5 1.2z m-29.5 17.1l18.4-18.5 6.3 6.3-18.4 18.4h-6.3v-6.2z"})})})}function d(e){let{editUrl:t}=e;return(0,i.jsxs)(a.A,{to:t,className:c.G.common.editThisPage,children:[(0,i.jsx)(l,{}),(0,i.jsx)(o.A,{id:"theme.common.editThisPage",description:"The link label to edit the current page",children:"Edit this page"})]})}var u=n(7638);function m(e){let{lastUpdatedAt:t}=e;const n=new Date(t),s=(0,u.i)({day:"numeric",month:"short",year:"numeric",timeZone:"UTC"}).format(n);return(0,i.jsx)(o.A,{id:"theme.lastUpdated.atDate",description:"The words used to describe on which date a page has been last updated",values:{date:(0,i.jsx)("b",{children:(0,i.jsx)("time",{dateTime:n.toISOString(),itemProp:"dateModified",children:s})})},children:" on {date}"})}function h(e){let{lastUpdatedBy:t}=e;return(0,i.jsx)(o.A,{id:"theme.lastUpdated.byUser",description:"The words used to describe by who the page has been last updated",values:{user:(0,i.jsx)("b",{children:t})},children:" by {user}"})}function p(e){let{lastUpdatedAt:t,lastUpdatedBy:n}=e;return(0,i.jsxs)("span",{className:c.G.common.lastUpdated,children:[(0,i.jsx)(o.A,{id:"theme.lastUpdated.lastUpdatedAtBy",description:"The sentence used to display when a page has been last updated, and by who",values:{atDate:t?(0,i.jsx)(m,{lastUpdatedAt:t}):"",byUser:n?(0,i.jsx)(h,{lastUpdatedBy:n}):""},children:"Last updated{atDate}{byUser}"}),!1]})}const f={lastUpdated:"lastUpdated_mIWg"};function x(e){let{className:t,editUrl:n,lastUpdatedAt:o,lastUpdatedBy:c}=e;return(0,i.jsxs)("div",{className:(0,s.A)("row",t),children:[(0,i.jsx)("div",{className:"col",children:n&&(0,i.jsx)(d,{editUrl:n})}),(0,i.jsx)("div",{className:(0,s.A)("col",f.lastUpdated),children:(o||c)&&(0,i.jsx)(p,{lastUpdatedAt:o,lastUpdatedBy:c})})]})}},3881:(e,t,n)=>{"use strict";n.d(t,{R:()=>a,x:()=>r});var s=n(8101);const o={},c=s.createContext(o);function a(e){const t=s.useContext(c);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function r(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:a(e.components),s.createElement(c.Provider,{value:t},e.children)}},4809:(e,t)=>{function n(e){let t,n=[];for(let s of e.split(",").map((e=>e.trim())))if(/^-?\d+$/.test(s))n.push(parseInt(s,10));else if(t=s.match(/^(-?\d+)(-|\.\.\.?|\u2025|\u2026|\u22EF)(-?\d+)$/)){let[e,s,o,c]=t;if(s&&c){s=parseInt(s),c=parseInt(c);const e=s<c?1:-1;"-"!==o&&".."!==o&&"\u2025"!==o||(c+=e);for(let t=s;t!==c;t+=e)n.push(t)}}return n}t.default=n,e.exports=n},7638:(e,t,n)=>{"use strict";n.d(t,{i:()=>o});var s=n(4291);function o(e){void 0===e&&(e={});const{i18n:{currentLocale:t}}=(0,s.A)(),n=function(){const{i18n:{currentLocale:e,localeConfigs:t}}=(0,s.A)();return t[e].calendar}();return new Intl.DateTimeFormat(t,{calendar:n,...e})}},7820:(e,t,n)=>{"use strict";n.d(t,{A:()=>U});var s=n(8101),o=n(5105);function c(e){const{mdxAdmonitionTitle:t,rest:n}=function(e){const t=s.Children.toArray(e),n=t.find((e=>s.isValidElement(e)&&"mdxAdmonitionTitle"===e.type)),c=t.filter((e=>e!==n)),a=n?.props.children;return{mdxAdmonitionTitle:a,rest:c.length>0?(0,o.jsx)(o.Fragment,{children:c}):null}}(e.children),c=e.title??t;return{...e,...c&&{title:c},children:n}}var a=n(3526),r=n(4323),i=n(6435);const l="admonition_kCc1",d="admonitionHeading_WWyw",u="admonitionIcon_tH_p",m="admonitionContent_lI_Q";function h(e){let{type:t,className:n,children:s}=e;return(0,o.jsx)("div",{className:(0,a.A)(i.G.common.admonition,i.G.common.admonitionType(t),l,n),children:s})}function p(e){let{icon:t,title:n}=e;return(0,o.jsxs)("div",{className:d,children:[(0,o.jsx)("span",{className:u,children:t}),n]})}function f(e){let{children:t}=e;return t?(0,o.jsx)("div",{className:m,children:t}):null}function x(e){const{type:t,icon:n,title:s,children:c,className:a}=e;return(0,o.jsxs)(h,{type:t,className:a,children:[s||n?(0,o.jsx)(p,{title:s,icon:n}):null,(0,o.jsx)(f,{children:c})]})}function j(e){return(0,o.jsx)("svg",{viewBox:"0 0 14 16",...e,children:(0,o.jsx)("path",{fillRule:"evenodd",d:"M6.3 5.69a.942.942 0 0 1-.28-.7c0-.28.09-.52.28-.7.19-.18.42-.28.7-.28.28 0 .52.09.7.28.18.19.28.42.28.7 0 .28-.09.52-.28.7a1 1 0 0 1-.7.3c-.28 0-.52-.11-.7-.3zM8 7.99c-.02-.25-.11-.48-.31-.69-.2-.19-.42-.3-.69-.31H6c-.27.02-.48.13-.69.31-.2.2-.3.44-.31.69h1v3c.02.27.11.5.31.69.2.2.42.31.69.31h1c.27 0 .48-.11.69-.31.2-.19.3-.42.31-.69H8V7.98v.01zM7 2.3c-3.14 0-5.7 2.54-5.7 5.68 0 3.14 2.56 5.7 5.7 5.7s5.7-2.55 5.7-5.7c0-3.15-2.56-5.69-5.7-5.69v.01zM7 .98c3.86 0 7 3.14 7 7s-3.14 7-7 7-7-3.12-7-7 3.14-7 7-7z"})})}const g={icon:(0,o.jsx)(j,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.note",description:"The default label used for the Note admonition (:::note)",children:"note"})};function b(e){return(0,o.jsx)(x,{...g,...e,className:(0,a.A)("alert alert--secondary",e.className),children:e.children})}function v(e){return(0,o.jsx)("svg",{viewBox:"0 0 12 16",...e,children:(0,o.jsx)("path",{fillRule:"evenodd",d:"M6.5 0C3.48 0 1 2.19 1 5c0 .92.55 2.25 1 3 1.34 2.25 1.78 2.78 2 4v1h5v-1c.22-1.22.66-1.75 2-4 .45-.75 1-2.08 1-3 0-2.81-2.48-5-5.5-5zm3.64 7.48c-.25.44-.47.8-.67 1.11-.86 1.41-1.25 2.06-1.45 3.23-.02.05-.02.11-.02.17H5c0-.06 0-.13-.02-.17-.2-1.17-.59-1.83-1.45-3.23-.2-.31-.42-.67-.67-1.11C2.44 6.78 2 5.65 2 5c0-2.2 2.02-4 4.5-4 1.22 0 2.36.42 3.22 1.19C10.55 2.94 11 3.94 11 5c0 .66-.44 1.78-.86 2.48zM4 14h5c-.23 1.14-1.3 2-2.5 2s-2.27-.86-2.5-2z"})})}const y={icon:(0,o.jsx)(v,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.tip",description:"The default label used for the Tip admonition (:::tip)",children:"tip"})};function N(e){return(0,o.jsx)(x,{...y,...e,className:(0,a.A)("alert alert--success",e.className),children:e.children})}function A(e){return(0,o.jsx)("svg",{viewBox:"0 0 14 16",...e,children:(0,o.jsx)("path",{fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"})})}const k={icon:(0,o.jsx)(A,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.info",description:"The default label used for the Info admonition (:::info)",children:"info"})};function B(e){return(0,o.jsx)(x,{...k,...e,className:(0,a.A)("alert alert--info",e.className),children:e.children})}function C(e){return(0,o.jsx)("svg",{viewBox:"0 0 16 16",...e,children:(0,o.jsx)("path",{fillRule:"evenodd",d:"M8.893 1.5c-.183-.31-.52-.5-.887-.5s-.703.19-.886.5L.138 13.499a.98.98 0 0 0 0 1.001c.193.31.53.501.886.501h13.964c.367 0 .704-.19.877-.5a1.03 1.03 0 0 0 .01-1.002L8.893 1.5zm.133 11.497H6.987v-2.003h2.039v2.003zm0-3.004H6.987V5.987h2.039v4.006z"})})}const w={icon:(0,o.jsx)(C,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.warning",description:"The default label used for the Warning admonition (:::warning)",children:"warning"})};function T(e){return(0,o.jsx)("svg",{viewBox:"0 0 12 16",...e,children:(0,o.jsx)("path",{fillRule:"evenodd",d:"M5.05.31c.81 2.17.41 3.38-.52 4.31C3.55 5.67 1.98 6.45.9 7.98c-1.45 2.05-1.7 6.53 3.53 7.7-2.2-1.16-2.67-4.52-.3-6.61-.61 2.03.53 3.33 1.94 2.86 1.39-.47 2.3.53 2.27 1.67-.02.78-.31 1.44-1.13 1.81 3.42-.59 4.78-3.42 4.78-5.56 0-2.84-2.53-3.22-1.25-5.61-1.52.13-2.03 1.13-1.89 2.75.09 1.08-1.02 1.8-1.86 1.33-.67-.41-.66-1.19-.06-1.78C8.18 5.31 8.68 2.45 5.05.32L5.03.3l.02.01z"})})}const E={icon:(0,o.jsx)(T,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.danger",description:"The default label used for the Danger admonition (:::danger)",children:"danger"})};const L={icon:(0,o.jsx)(C,{}),title:(0,o.jsx)(r.A,{id:"theme.admonition.caution",description:"The default label used for the Caution admonition (:::caution)",children:"caution"})};const _={...{note:b,tip:N,info:B,warning:function(e){return(0,o.jsx)(x,{...w,...e,className:(0,a.A)("alert alert--warning",e.className),children:e.children})},danger:function(e){return(0,o.jsx)(x,{...E,...e,className:(0,a.A)("alert alert--danger",e.className),children:e.children})}},...{secondary:e=>(0,o.jsx)(b,{title:"secondary",...e}),important:e=>(0,o.jsx)(B,{title:"important",...e}),success:e=>(0,o.jsx)(N,{title:"success",...e}),caution:function(e){return(0,o.jsx)(x,{...L,...e,className:(0,a.A)("alert alert--warning",e.className),children:e.children})}}};function U(e){const t=c(e),n=(s=t.type,_[s]||(console.warn(`No admonition component found for admonition type "${s}". Using Info as fallback.`),_.info));var s;return(0,o.jsx)(n,{...t})}},9131:(e,t,n)=>{"use strict";n.d(t,{A:()=>me});var s=n(8101),o=n(3881),c=n(6017),a=n(634),r=n(3526),i=n(4761),l=n(4538);function d(){const{prism:e}=(0,l.p)(),{colorMode:t}=(0,i.G)(),n=e.theme,s=e.darkTheme||n;return"dark"===t?s:n}var u=n(6435),m=n(4809),h=n.n(m);const p=/title=(?<quote>["'])(?<title>.*?)\1/,f=/\{(?<range>[\d,-]+)\}/,x={js:{start:"\\/\\/",end:""},jsBlock:{start:"\\/\\*",end:"\\*\\/"},jsx:{start:"\\{\\s*\\/\\*",end:"\\*\\/\\s*\\}"},bash:{start:"#",end:""},html:{start:"\x3c!--",end:"--\x3e"}},j={...x,lua:{start:"--",end:""},wasm:{start:"\\;\\;",end:""},tex:{start:"%",end:""},vb:{start:"['\u2018\u2019]",end:""},vbnet:{start:"(?:_\\s*)?['\u2018\u2019]",end:""},rem:{start:"[Rr][Ee][Mm]\\b",end:""},f90:{start:"!",end:""},ml:{start:"\\(\\*",end:"\\*\\)"},cobol:{start:"\\*>",end:""}},g=Object.keys(x);function b(e,t){const n=e.map((e=>{const{start:n,end:s}=j[e];return`(?:${n}\\s*(${t.flatMap((e=>[e.line,e.block?.start,e.block?.end].filter(Boolean))).join("|")})\\s*${s})`})).join("|");return new RegExp(`^\\s*(?:${n})\\s*$`)}function v(e,t){let n=e.replace(/\n$/,"");const{language:s,magicComments:o,metastring:c}=t;if(c&&f.test(c)){const e=c.match(f).groups.range;if(0===o.length)throw new Error(`A highlight range has been given in code block's metastring (\`\`\` ${c}), but no magic comment config is available. Docusaurus applies the first magic comment entry's className for metastring ranges.`);const t=o[0].className,s=h()(e).filter((e=>e>0)).map((e=>[e-1,[t]]));return{lineClassNames:Object.fromEntries(s),code:n}}if(void 0===s)return{lineClassNames:{},code:n};const a=function(e,t){switch(e){case"js":case"javascript":case"ts":case"typescript":return b(["js","jsBlock"],t);case"jsx":case"tsx":return b(["js","jsBlock","jsx"],t);case"html":return b(["js","jsBlock","html"],t);case"python":case"py":case"bash":return b(["bash"],t);case"markdown":case"md":return b(["html","jsx","bash"],t);case"tex":case"latex":case"matlab":return b(["tex"],t);case"lua":case"haskell":return b(["lua"],t);case"sql":return b(["lua","jsBlock"],t);case"wasm":return b(["wasm"],t);case"vb":case"vba":case"visual-basic":return b(["vb","rem"],t);case"vbnet":return b(["vbnet","rem"],t);case"batch":return b(["rem"],t);case"basic":return b(["rem","f90"],t);case"fsharp":return b(["js","ml"],t);case"ocaml":case"sml":return b(["ml"],t);case"fortran":return b(["f90"],t);case"cobol":return b(["cobol"],t);default:return b(g,t)}}(s,o),r=n.split("\n"),i=Object.fromEntries(o.map((e=>[e.className,{start:0,range:""}]))),l=Object.fromEntries(o.filter((e=>e.line)).map((e=>{let{className:t,line:n}=e;return[n,t]}))),d=Object.fromEntries(o.filter((e=>e.block)).map((e=>{let{className:t,block:n}=e;return[n.start,t]}))),u=Object.fromEntries(o.filter((e=>e.block)).map((e=>{let{className:t,block:n}=e;return[n.end,t]})));for(let h=0;h<r.length;){const e=r[h].match(a);if(!e){h+=1;continue}const t=e.slice(1).find((e=>void 0!==e));l[t]?i[l[t]].range+=`${h},`:d[t]?i[d[t]].start=h:u[t]&&(i[u[t]].range+=`${i[u[t]].start}-${h-1},`),r.splice(h,1)}n=r.join("\n");const m={};return Object.entries(i).forEach((e=>{let[t,{range:n}]=e;h()(n).forEach((e=>{m[e]??=[],m[e].push(t)}))})),{lineClassNames:m,code:n}}const y="codeBlockContainer_HTPM";var N=n(5105);function A(e){let{as:t,...n}=e;const s=function(e){const t={color:"--prism-color",backgroundColor:"--prism-background-color"},n={};return Object.entries(e.plain).forEach((e=>{let[s,o]=e;const c=t[s];c&&"string"==typeof o&&(n[c]=o)})),n}(d());return(0,N.jsx)(t,{...n,style:s,className:(0,r.A)(n.className,y,u.G.common.codeBlock)})}const k={codeBlockContent:"codeBlockContent_OwhC",codeBlockTitle:"codeBlockTitle_OQJp",codeBlock:"codeBlock_QQmy",codeBlockStandalone:"codeBlockStandalone_zb5l",codeBlockLines:"codeBlockLines_piqd",codeBlockLinesWithNumbering:"codeBlockLinesWithNumbering_u83U",buttonGroup:"buttonGroup_n6sy"};function B(e){let{children:t,className:n}=e;return(0,N.jsx)(A,{as:"pre",tabIndex:0,className:(0,r.A)(k.codeBlockStandalone,"thin-scrollbar",n),children:(0,N.jsx)("code",{className:k.codeBlockLines,children:t})})}var C=n(5544);const w={attributes:!0,characterData:!0,childList:!0,subtree:!0};function T(e,t){const[n,o]=(0,s.useState)(),c=(0,s.useCallback)((()=>{o(e.current?.closest("[role=tabpanel][hidden]"))}),[e,o]);(0,s.useEffect)((()=>{c()}),[c]),function(e,t,n){void 0===n&&(n=w);const o=(0,C._q)(t),c=(0,C.Be)(n);(0,s.useEffect)((()=>{const t=new MutationObserver(o);return e&&t.observe(e,c),()=>t.disconnect()}),[e,o,c])}(n,(e=>{e.forEach((e=>{"attributes"===e.type&&"hidden"===e.attributeName&&(t(),c())}))}),{attributes:!0,characterData:!1,childList:!1,subtree:!1})}var E=n(2355);const L="codeLine_H69Z",_="codeLineNumber_m1Ph",U="codeLineContent_tv57";function M(e){let{line:t,classNames:n,showLineNumbers:s,getLineProps:o,getTokenProps:c}=e;1===t.length&&"\n"===t[0].content&&(t[0].content="");const a=o({line:t,className:(0,r.A)(n,s&&L)}),i=t.map(((e,t)=>(0,N.jsx)("span",{...c({token:e})},t)));return(0,N.jsxs)("span",{...a,children:[s?(0,N.jsxs)(N.Fragment,{children:[(0,N.jsx)("span",{className:_}),(0,N.jsx)("span",{className:U,children:i})]}):i,(0,N.jsx)("br",{})]})}var S=n(4323);function z(e){return(0,N.jsx)("svg",{viewBox:"0 0 24 24",...e,children:(0,N.jsx)("path",{fill:"currentColor",d:"M19,21H8V7H19M19,5H8A2,2 0 0,0 6,7V21A2,2 0 0,0 8,23H19A2,2 0 0,0 21,21V7A2,2 0 0,0 19,5M16,1H4A2,2 0 0,0 2,3V17H4V3H16V1Z"})})}function I(e){return(0,N.jsx)("svg",{viewBox:"0 0 24 24",...e,children:(0,N.jsx)("path",{fill:"currentColor",d:"M21,7L9,19L3.5,13.5L4.91,12.09L9,16.17L19.59,5.59L21,7Z"})})}const H={copyButtonCopied:"copyButtonCopied_lXZu",copyButtonIcons:"copyButtonIcons_GuCM",copyButtonIcon:"copyButtonIcon_KWtG",copyButtonSuccessIcon:"copyButtonSuccessIcon_nnIb"};function R(e){let{code:t,className:n}=e;const[o,c]=(0,s.useState)(!1),a=(0,s.useRef)(void 0),i=(0,s.useCallback)((()=>{!function(e,t){let{target:n=document.body}=void 0===t?{}:t;if("string"!=typeof e)throw new TypeError(`Expected parameter \`text\` to be a \`string\`, got \`${typeof e}\`.`);const s=document.createElement("textarea"),o=document.activeElement;s.value=e,s.setAttribute("readonly",""),s.style.contain="strict",s.style.position="absolute",s.style.left="-9999px",s.style.fontSize="12pt";const c=document.getSelection(),a=c.rangeCount>0&&c.getRangeAt(0);n.append(s),s.select(),s.selectionStart=0,s.selectionEnd=e.length;let r=!1;try{r=document.execCommand("copy")}catch{}s.remove(),a&&(c.removeAllRanges(),c.addRange(a)),o&&o.focus()}(t),c(!0),a.current=window.setTimeout((()=>{c(!1)}),1e3)}),[t]);return(0,s.useEffect)((()=>()=>window.clearTimeout(a.current)),[]),(0,N.jsx)("button",{type:"button","aria-label":o?(0,S.T)({id:"theme.CodeBlock.copied",message:"Copied",description:"The copied button label on code blocks"}):(0,S.T)({id:"theme.CodeBlock.copyButtonAriaLabel",message:"Copy code to clipboard",description:"The ARIA label for copy code blocks button"}),title:(0,S.T)({id:"theme.CodeBlock.copy",message:"Copy",description:"The copy button label on code blocks"}),className:(0,r.A)("clean-btn",n,H.copyButton,o&&H.copyButtonCopied),onClick:i,children:(0,N.jsxs)("span",{className:H.copyButtonIcons,"aria-hidden":"true",children:[(0,N.jsx)(z,{className:H.copyButtonIcon}),(0,N.jsx)(I,{className:H.copyButtonSuccessIcon})]})})}function V(e){return(0,N.jsx)("svg",{viewBox:"0 0 24 24",...e,children:(0,N.jsx)("path",{fill:"currentColor",d:"M4 19h6v-2H4v2zM20 5H4v2h16V5zm-3 6H4v2h13.25c1.1 0 2 .9 2 2s-.9 2-2 2H15v-2l-3 3l3 3v-2h2c2.21 0 4-1.79 4-4s-1.79-4-4-4z"})})}const P="wordWrapButtonIcon_jKCa",$="wordWrapButtonEnabled_CAa8";function W(e){let{className:t,onClick:n,isEnabled:s}=e;const o=(0,S.T)({id:"theme.CodeBlock.wordWrapToggle",message:"Toggle word wrap",description:"The title attribute for toggle word wrapping button of code block lines"});return(0,N.jsx)("button",{type:"button",onClick:n,className:(0,r.A)("clean-btn",t,s&&$),"aria-label":o,title:o,children:(0,N.jsx)(V,{className:P,"aria-hidden":"true"})})}function D(e){let{children:t,className:n="",metastring:o,title:c,showLineNumbers:a,language:i}=e;const{prism:{defaultLanguage:u,magicComments:m}}=(0,l.p)(),h=function(e){return e?.toLowerCase()}(i??function(e){const t=e.split(" ").find((e=>e.startsWith("language-")));return t?.replace(/language-/,"")}(n)??u),f=d(),x=function(){const[e,t]=(0,s.useState)(!1),[n,o]=(0,s.useState)(!1),c=(0,s.useRef)(null),a=(0,s.useCallback)((()=>{const n=c.current.querySelector("code");e?n.removeAttribute("style"):(n.style.whiteSpace="pre-wrap",n.style.overflowWrap="anywhere"),t((e=>!e))}),[c,e]),r=(0,s.useCallback)((()=>{const{scrollWidth:e,clientWidth:t}=c.current,n=e>t||c.current.querySelector("code").hasAttribute("style");o(n)}),[c]);return T(c,r),(0,s.useEffect)((()=>{r()}),[e,r]),(0,s.useEffect)((()=>(window.addEventListener("resize",r,{passive:!0}),()=>{window.removeEventListener("resize",r)})),[r]),{codeBlockRef:c,isEnabled:e,isCodeScrollable:n,toggle:a}}(),j=function(e){return e?.match(p)?.groups.title??""}(o)||c,{lineClassNames:g,code:b}=v(t,{metastring:o,language:h,magicComments:m}),y=a??function(e){return Boolean(e?.includes("showLineNumbers"))}(o);return(0,N.jsxs)(A,{as:"div",className:(0,r.A)(n,h&&!n.includes(`language-${h}`)&&`language-${h}`),children:[j&&(0,N.jsx)("div",{className:k.codeBlockTitle,children:j}),(0,N.jsxs)("div",{className:k.codeBlockContent,children:[(0,N.jsx)(E.f4,{theme:f,code:b,language:h??"text",children:e=>{let{className:t,style:n,tokens:s,getLineProps:o,getTokenProps:c}=e;return(0,N.jsx)("pre",{tabIndex:0,ref:x.codeBlockRef,className:(0,r.A)(t,k.codeBlock,"thin-scrollbar"),style:n,children:(0,N.jsx)("code",{className:(0,r.A)(k.codeBlockLines,y&&k.codeBlockLinesWithNumbering),children:s.map(((e,t)=>(0,N.jsx)(M,{line:e,getLineProps:o,getTokenProps:c,classNames:g[t],showLineNumbers:y},t)))})})}}),(0,N.jsxs)("div",{className:k.buttonGroup,children:[(x.isEnabled||x.isCodeScrollable)&&(0,N.jsx)(W,{className:k.codeButton,onClick:()=>x.toggle(),isEnabled:x.isEnabled}),(0,N.jsx)(R,{className:k.codeButton,code:b})]})]})]})}function G(e){let{children:t,...n}=e;const o=(0,a.A)(),c=function(e){return s.Children.toArray(e).some((e=>(0,s.isValidElement)(e)))?e:Array.isArray(e)?e.join(""):e}(t),r="string"==typeof c?D:B;return(0,N.jsx)(r,{...n,children:c},String(o))}function O(e){return(0,N.jsx)("code",{...e})}var q=n(9515);var F=n(3139),Z=n(8714);const Q="details_BChi",K="isBrowser_yjV1",J="collapsibleContent_fnwV";function X(e){return!!e&&("SUMMARY"===e.tagName||X(e.parentElement))}function Y(e,t){return!!e&&(e===t||Y(e.parentElement,t))}function ee(e){let{summary:t,children:n,...o}=e;(0,F.A)().collectAnchor(o.id);const c=(0,a.A)(),i=(0,s.useRef)(null),{collapsed:l,setCollapsed:d}=(0,Z.u)({initialState:!o.open}),[u,m]=(0,s.useState)(o.open),h=s.isValidElement(t)?t:(0,N.jsx)("summary",{children:t??"Details"});return(0,N.jsxs)("details",{...o,ref:i,open:u,"data-collapsed":l,className:(0,r.A)(Q,c&&K,o.className),onMouseDown:e=>{X(e.target)&&e.detail>1&&e.preventDefault()},onClick:e=>{e.stopPropagation();const t=e.target;X(t)&&Y(t,i.current)&&(e.preventDefault(),l?(d(!1),m(!0)):d(!0))},children:[h,(0,N.jsx)(Z.N,{lazy:!1,collapsed:l,disableSSRStyle:!0,onCollapseTransitionEnd:e=>{d(e),m(!e)},children:(0,N.jsx)("div",{className:J,children:n})})]})}const te="details_TCPE";function ne(e){let{...t}=e;return(0,N.jsx)(ee,{...t,className:(0,r.A)("alert alert--info",te,t.className)})}function se(e){const t=s.Children.toArray(e.children),n=t.find((e=>s.isValidElement(e)&&"summary"===e.type)),o=(0,N.jsx)(N.Fragment,{children:t.filter((e=>e!==n))});return(0,N.jsx)(ne,{...e,summary:n,children:o})}var oe=n(2068);function ce(e){return(0,N.jsx)(oe.A,{...e})}const ae="containsTaskList_mIPT";function re(e){if(void 0!==e)return(0,r.A)(e,e?.includes("contains-task-list")&&ae)}const ie="img__G2c";var le=n(7820),de=n(5099);const ue={Head:c.A,details:se,Details:se,code:function(e){return function(e){return void 0!==e.children&&s.Children.toArray(e.children).every((e=>"string"==typeof e&&!e.includes("\n")))}(e)?(0,N.jsx)(O,{...e}):(0,N.jsx)(G,{...e})},a:function(e){return(0,N.jsx)(q.A,{...e})},pre:function(e){return(0,N.jsx)(N.Fragment,{children:e.children})},ul:function(e){return(0,N.jsx)("ul",{...e,className:re(e.className)})},li:function(e){return(0,F.A)().collectAnchor(e.id),(0,N.jsx)("li",{...e})},img:function(e){return(0,N.jsx)("img",{decoding:"async",loading:"lazy",...e,className:(t=e.className,(0,r.A)(t,ie))});var t},h1:e=>(0,N.jsx)(ce,{as:"h1",...e}),h2:e=>(0,N.jsx)(ce,{as:"h2",...e}),h3:e=>(0,N.jsx)(ce,{as:"h3",...e}),h4:e=>(0,N.jsx)(ce,{as:"h4",...e}),h5:e=>(0,N.jsx)(ce,{as:"h5",...e}),h6:e=>(0,N.jsx)(ce,{as:"h6",...e}),admonition:le.A,mermaid:de.A};function me(e){let{children:t}=e;return(0,N.jsx)(o.x,{components:ue,children:t})}}}]);