"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7723],{3673:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>c,contentTitle:()=>a,default:()=>h,frontMatter:()=>o,metadata:()=>i,toc:()=>l});const i=JSON.parse('{"id":"web-ui","title":"Web UI Configuration","description":"ServiceRadar includes a modern web interface built with Next.js that provides a dashboard for monitoring your infrastructure. This guide explains how to install, configure, and secure the web UI component.","source":"@site/docs/web-ui.md","sourceDirName":".","slug":"/web-ui","permalink":"/docs/web-ui","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":5,"frontMatter":{"sidebar_position":5,"title":"Web UI Configuration"},"sidebar":"tutorialSidebar","previous":{"title":"TLS Security","permalink":"/docs/tls-security"},"next":{"title":"Architecture","permalink":"/docs/architecture"}}');var s=r(5105),t=r(3881);const o={sidebar_position:5,title:"Web UI Configuration"},a="Web UI Configuration",c={},l=[{value:"Overview",id:"overview",level:2},{value:"Architecture",id:"architecture",level:2},{value:"Installation",id:"installation",level:2},{value:"Configuration",id:"configuration",level:2},{value:"Web UI Configuration",id:"web-ui-configuration-1",level:3},{value:"Nginx Configuration",id:"nginx-configuration",level:3},{value:"API Key Security",id:"api-key-security",level:2},{value:"Security Features",id:"security-features",level:2},{value:"Custom Domain and SSL",id:"custom-domain-and-ssl",level:2},{value:"Troubleshooting",id:"troubleshooting",level:2}];function d(e){const n={admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",mermaid:"mermaid",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,t.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.header,{children:(0,s.jsx)(n.h1,{id:"web-ui-configuration",children:"Web UI Configuration"})}),"\n",(0,s.jsx)(n.p,{children:"ServiceRadar includes a modern web interface built with Next.js that provides a dashboard for monitoring your infrastructure. This guide explains how to install, configure, and secure the web UI component."}),"\n",(0,s.jsx)(n.h2,{id:"overview",children:"Overview"}),"\n",(0,s.jsx)(n.p,{children:"The ServiceRadar web interface:"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"Provides a visual dashboard for monitoring your infrastructure"}),"\n",(0,s.jsx)(n.li,{children:"Communicates securely with the ServiceRadar API"}),"\n",(0,s.jsx)(n.li,{children:"Uses Nginx as a reverse proxy to handle HTTP requests"}),"\n",(0,s.jsx)(n.li,{children:"Automatically configures security with API key authentication"}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"architecture",children:"Architecture"}),"\n",(0,s.jsx)(n.mermaid,{value:'graph LR\n    subgraph "ServiceRadar Server"\n        A[Web Browser] --\x3e|HTTP/80| B[Nginx]\n        B --\x3e|Proxy /api| C[Core API<br/>:8090]\n        B --\x3e|Proxy /| D[Next.js<br/>:3000]\n        D --\x3e|API Requests with<br/>Injected API Key| C\n    end'}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Nginx"})," runs on port 80 and acts as the main entry point"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Next.js"})," provides the web UI on port 3000"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Core API"})," service runs on port 8090"]}),"\n",(0,s.jsx)(n.li,{children:"API requests from the UI are secured with an automatically generated API key"}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"installation",children:"Installation"}),"\n",(0,s.jsxs)(n.p,{children:["The web UI is installed via the ",(0,s.jsx)(n.code,{children:"serviceradar-web"})," package:"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"curl -LO https://github.com/carverauto/serviceradar/releases/download/1.0.21/serviceradar-web_1.0.21.deb\nsudo dpkg -i serviceradar-web_1.0.21.deb\n"})}),"\n",(0,s.jsx)(n.admonition,{type:"note",children:(0,s.jsxs)(n.p,{children:["It's recommended to install the ",(0,s.jsx)(n.code,{children:"serviceradar-core"})," package first, as the web UI depends on it."]})}),"\n",(0,s.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,s.jsx)(n.h3,{id:"web-ui-configuration-1",children:"Web UI Configuration"}),"\n",(0,s.jsxs)(n.p,{children:["Edit ",(0,s.jsx)(n.code,{children:"/etc/serviceradar/web.json"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-json",children:'{\n  "port": 3000,\n  "host": "0.0.0.0",\n  "api_url": "http://localhost:8090"\n}\n'})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.code,{children:"port"}),": The port for the Next.js application (default: 3000)"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.code,{children:"host"}),": The host address to bind to"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.code,{children:"api_url"}),": The URL for the core API service"]}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"nginx-configuration",children:"Nginx Configuration"}),"\n",(0,s.jsxs)(n.p,{children:["The web package automatically configures Nginx. The main configuration file is located at ",(0,s.jsx)(n.code,{children:"/etc/nginx/conf.d/serviceradar-web.conf"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-nginx",children:'# ServiceRadar Web Interface - Nginx Configuration\nserver {\n    listen 80;\n    server_name _; # Catch-all server name (use your domain if you have one)\n\n    access_log /var/log/nginx/serviceradar-web.access.log;\n    error_log /var/log/nginx/serviceradar-web.error.log;\n\n    # API proxy (assumes serviceradar-core package is installed)\n    location /api/ {\n        proxy_pass http://localhost:8090;\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n        proxy_set_header X-Forwarded-Proto $scheme;\n    }\n\n    # Support for Next.js WebSockets (if used)\n    location /_next/webpack-hmr {\n        proxy_pass http://localhost:3000;\n        proxy_http_version 1.1;\n        proxy_set_header Upgrade $http_upgrade;\n        proxy_set_header Connection "upgrade";\n    }\n\n    # Main app - proxy all requests to Next.js\n    location / {\n        proxy_pass http://127.0.0.1:3000;\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n        proxy_set_header X-Forwarded-Proto $scheme;\n    }\n}\n'})}),"\n",(0,s.jsx)(n.p,{children:"You can customize this file for your specific domain or add SSL configuration."}),"\n",(0,s.jsx)(n.h2,{id:"api-key-security",children:"API Key Security"}),"\n",(0,s.jsx)(n.p,{children:"ServiceRadar uses an API key to secure communication between the web UI and the core API. This key is automatically generated during installation and stored in:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{children:"/etc/serviceradar/api.env\n"})}),"\n",(0,s.jsx)(n.p,{children:"The file contains an environment variable:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{children:"API_KEY=your_generated_key\n"})}),"\n",(0,s.jsx)(n.p,{children:"This API key is:"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["Automatically generated during ",(0,s.jsx)(n.code,{children:"serviceradar-core"})," installation"]}),"\n",(0,s.jsx)(n.li,{children:"Used by the web UI's Next.js middleware to authenticate API requests"}),"\n",(0,s.jsx)(n.li,{children:"Securely injected into backend requests without exposing it to clients"}),"\n"]}),"\n",(0,s.jsx)(n.admonition,{type:"caution",children:(0,s.jsxs)(n.p,{children:["Keep the API key secure. Don't share the content of the ",(0,s.jsx)(n.code,{children:"api.env"})," file."]})}),"\n",(0,s.jsx)(n.h2,{id:"security-features",children:"Security Features"}),"\n",(0,s.jsx)(n.p,{children:"The web UI includes several security features:"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Server-side rendering"}),": The Next.js application runs in SSR mode, which keeps sensitive code on the server"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"API middleware"}),": Requests to the Core API are handled via middleware that injects the API key"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Proxy architecture"}),": Direct client access to the API is prevented through the proxy setup"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"Isolation"}),": The web UI runs as a separate service with limited permissions"]}),"\n"]}),"\n",(0,s.jsx)(n.h2,{id:"custom-domain-and-ssl",children:"Custom Domain and SSL"}),"\n",(0,s.jsx)(n.p,{children:"To configure a custom domain with SSL:"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsx)(n.li,{children:"Update the Nginx configuration with your domain name"}),"\n",(0,s.jsx)(n.li,{children:"Add SSL certificate configuration"}),"\n",(0,s.jsx)(n.li,{children:"Restart Nginx"}),"\n"]}),"\n",(0,s.jsx)(n.p,{children:"Example configuration with SSL:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-nginx",children:"server {\n    listen 80;\n    server_name your-domain.com;\n    return 301 https://$host$request_uri;\n}\n\nserver {\n    listen 443 ssl;\n    server_name your-domain.com;\n\n    ssl_certificate /path/to/your/certificate.crt;\n    ssl_certificate_key /path/to/your/private.key;\n    \n    # ... rest of the configuration\n}\n"})}),"\n",(0,s.jsx)(n.h2,{id:"troubleshooting",children:"Troubleshooting"}),"\n",(0,s.jsx)(n.p,{children:"Common issues and solutions:"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Web UI not accessible"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:["Check if Nginx is running: ",(0,s.jsx)(n.code,{children:"systemctl status nginx"})]}),"\n",(0,s.jsxs)(n.li,{children:["Verify the Next.js application is running: ",(0,s.jsx)(n.code,{children:"systemctl status serviceradar-web"})]}),"\n",(0,s.jsxs)(n.li,{children:["Check ports: ",(0,s.jsx)(n.code,{children:"netstat -tulpn | grep -E '3000|80'"})]}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"API connection errors"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:["Verify the Core API is running: ",(0,s.jsx)(n.code,{children:"systemctl status serviceradar-core"})]}),"\n",(0,s.jsx)(n.li,{children:"Check API key exists and is properly formatted"}),"\n",(0,s.jsx)(n.li,{children:"Verify API URL in web.json is correct"}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Permission issues"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:["Check ownership of files: ",(0,s.jsx)(n.code,{children:"ls -la /etc/serviceradar/"})]}),"\n",(0,s.jsx)(n.li,{children:"Ensure the serviceradar user has appropriate permissions"}),"\n"]}),"\n"]}),"\n",(0,s.jsxs)(n.li,{children:["\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"Nginx configuration errors"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:["Test configuration: ",(0,s.jsx)(n.code,{children:"nginx -t"})]}),"\n",(0,s.jsxs)(n.li,{children:["Check logs: ",(0,s.jsx)(n.code,{children:"tail -f /var/log/nginx/error.log"})]}),"\n"]}),"\n"]}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,t.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},3881:(e,n,r)=>{r.d(n,{R:()=>o,x:()=>a});var i=r(8101);const s={},t=i.createContext(s);function o(e){const n=i.useContext(t);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:o(e.components),i.createElement(t.Provider,{value:n},e.children)}}}]);