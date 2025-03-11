"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[3873],{3757:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>l,contentTitle:()=>t,default:()=>h,frontMatter:()=>c,metadata:()=>s,toc:()=>d});const s=JSON.parse('{"id":"configuration","title":"Configuration Basics","description":"ServiceRadar components are configured via JSON files in /etc/serviceradar/. This guide covers the essential configurations needed to get your monitoring system up and running.","source":"@site/docs/configuration.md","sourceDirName":".","slug":"/configuration","permalink":"/docs/configuration","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":3,"frontMatter":{"sidebar_position":3,"title":"Configuration Basics"},"sidebar":"tutorialSidebar","previous":{"title":"Installation Guide","permalink":"/docs/installation"},"next":{"title":"TLS Security","permalink":"/docs/tls-security"}}');var i=r(5105),o=r(3881);const c={sidebar_position:3,title:"Configuration Basics"},t="Configuration Basics",l={},d=[{value:"Agent Configuration",id:"agent-configuration",level:2},{value:"Configuration Options:",id:"configuration-options",level:3},{value:"Poller Configuration",id:"poller-configuration",level:2},{value:"Configuration Options:",id:"configuration-options-1",level:3},{value:"Check Types:",id:"check-types",level:3},{value:"Core Configuration",id:"core-configuration",level:2},{value:"API Key",id:"api-key",level:3},{value:"Dusk Node Checker",id:"dusk-node-checker",level:3},{value:"Network Sweep",id:"network-sweep",level:3},{value:"Next Steps",id:"next-steps",level:2}];function a(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",ul:"ul",...(0,o.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.header,{children:(0,i.jsx)(n.h1,{id:"configuration-basics",children:"Configuration Basics"})}),"\n",(0,i.jsxs)(n.p,{children:["ServiceRadar components are configured via JSON files in ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/"}),". This guide covers the essential configurations needed to get your monitoring system up and running."]}),"\n",(0,i.jsx)(n.h2,{id:"agent-configuration",children:"Agent Configuration"}),"\n",(0,i.jsx)(n.p,{children:"The agent runs on each monitored host and collects status information from services."}),"\n",(0,i.jsxs)(n.p,{children:["Edit ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/agent.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "checkers_dir": "/etc/serviceradar/checkers",\n  "listen_addr": ":50051",\n  "service_type": "grpc",\n  "service_name": "AgentService",\n  "security": {\n    "mode": "none",\n    "cert_dir": "/etc/serviceradar/certs",\n    "server_name": "changeme",\n    "role": "agent"\n  }\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"configuration-options",children:"Configuration Options:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"checkers_dir"}),": Directory containing checker configurations"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"listen_addr"}),": Address and port the agent listens on"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"service_type"}),': Type of service (should be "grpc")']}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"security"}),": Security settings","\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"mode"}),': Security mode ("none" or "mtls")']}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"cert_dir"}),": Directory for TLS certificates"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"server_name"}),": Hostname/IP of the poller (important for TLS)"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"role"}),': Role of this component ("agent")']}),"\n"]}),"\n"]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"poller-configuration",children:"Poller Configuration"}),"\n",(0,i.jsx)(n.p,{children:"The poller contacts agents to collect monitoring data and reports to the core service."}),"\n",(0,i.jsxs)(n.p,{children:["Edit ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/poller.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "agents": {\n    "local-agent": {\n      "address": "localhost:50051",\n      "security": { \n        "server_name": "changeme", \n        "mode": "none" \n      },\n      "checks": [\n        { "service_type": "process", "service_name": "nginx", "details": "nginx" },\n        { "service_type": "port", "service_name": "SSH", "details": "127.0.0.1:22" },\n        { "service_type": "icmp", "service_name": "ping", "details": "8.8.8.8" }\n      ]\n    }\n  },\n  "core_address": "changeme:50052",\n  "listen_addr": ":50053",\n  "poll_interval": "30s",\n  "poller_id": "my-poller",\n  "service_name": "PollerService",\n  "service_type": "grpc",\n  "security": {\n    "mode": "none",\n    "cert_dir": "/etc/serviceradar/certs",\n    "server_name": "changeme",\n    "role": "poller"\n  }\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"configuration-options-1",children:"Configuration Options:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"agents"}),": Map of agents to monitor","\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Each agent has an ",(0,i.jsx)(n.code,{children:"address"}),", ",(0,i.jsx)(n.code,{children:"security"})," settings, and ",(0,i.jsx)(n.code,{children:"checks"})," to perform"]}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"core_address"}),": Address of the core service"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"listen_addr"}),": Address and port the poller listens on"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"poll_interval"}),": How often to poll agents"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"poller_id"}),": Unique identifier for this poller"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"security"}),": Security settings (similar to agent)"]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"check-types",children:"Check Types:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"process"}),": Check if a process is running"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"port"}),": Check if a TCP port is responding"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"icmp"}),": Ping a host"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"grpc"}),": Check a gRPC service"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"snmp"}),": Check via SNMP (requires snmp checker)"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"sweep"}),": Network sweep check"]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"core-configuration",children:"Core Configuration"}),"\n",(0,i.jsx)(n.p,{children:"The core service receives reports from pollers and provides the API backend."}),"\n",(0,i.jsxs)(n.p,{children:["Edit ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/core.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "listen_addr": ":8090",\n  "grpc_addr": ":50052",\n  "alert_threshold": "5m",\n  "known_pollers": ["my-poller"],\n  "metrics": {\n    "enabled": true,\n    "retention": 100,\n    "max_nodes": 10000\n  },\n  "security": {\n    "mode": "none",\n    "cert_dir": "/etc/serviceradar/certs",\n    "role": "core"\n  },\n  "webhooks": [\n    {\n      "enabled": false,\n      "url": "https://your-webhook-url",\n      "cooldown": "15m",\n      "headers": [\n        {\n          "key": "Authorization",\n          "value": "Bearer your-token"\n        }\n      ]\n    },\n    {\n      "enabled": true,\n      "url": "https://discord.com/api/webhooks/changeme",\n      "cooldown": "15m",\n      "template": "{\\"embeds\\":[{\\"title\\":\\"{{.alert.Title}}\\",\\"description\\":\\"{{.alert.Message}}\\",\\"color\\":{{if eq .alert.Level \\"error\\"}}15158332{{else if eq .alert.Level \\"warning\\"}}16776960{{else}}3447003{{end}},\\"timestamp\\":\\"{{.alert.Timestamp}}\\",\\"fields\\":[{\\"name\\":\\"Node ID\\",\\"value\\":\\"{{.alert.NodeID}}\\",\\"inline\\":true}{{range $key, $value := .alert.Details}},{\\"name\\":\\"{{$key}}\\",\\"value\\":\\"{{$value}}\\",\\"inline\\":true}{{end}}]}]}"\n    }\n  ]\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"api-key",children:"API Key"}),"\n",(0,i.jsx)(n.p,{children:"During installation, the core service automatically generates an API key, stored in:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{children:"/etc/serviceradar/api.env\n"})}),"\n",(0,i.jsx)(n.p,{children:"This API key is used for secure communication between the web UI and the core API. The key is automatically injected into API requests by the web UI's middleware, ensuring secure communication without exposing the key to clients."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{children:'\n### Configuration Options:\n\n- `listen_addr`: Address and port for web dashboard\n- `grpc_addr`: Address and port for gRPC service\n- `alert_threshold`: How long a service must be down before alerting\n- `known_pollers`: List of poller IDs that can connect\n- `metrics`: Metrics collection settings\n- `security`: Security settings (similar to agent)\n- `webhooks`: List of webhook configurations for alerts\n\n## Optional Checker Configurations\n\n### SNMP Checker\n\nFor monitoring network devices via SNMP, edit `/etc/serviceradar/checkers/snmp.json`:\n\n```json\n{\n  "node_address": "localhost:50051",\n  "listen_addr": ":50054",\n  "security": {\n    "server_name": "changeme",\n    "mode": "none",\n    "role": "checker",\n    "cert_dir": "/etc/serviceradar/certs"\n  },\n  "timeout": "30s",\n  "targets": [\n    {\n      "name": "router",\n      "host": "192.168.1.1",\n      "port": 161,\n      "community": "public",\n      "version": "v2c",\n      "interval": "30s",\n      "retries": 2,\n      "oids": [\n        {\n          "oid": ".1.3.6.1.2.1.2.2.1.10.4",\n          "name": "ifInOctets_4",\n          "type": "counter",\n          "scale": 1.0\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"dusk-node-checker",children:"Dusk Node Checker"}),"\n",(0,i.jsxs)(n.p,{children:["For monitoring Dusk nodes, edit ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/checkers/dusk.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "name": "dusk",\n  "type": "grpc",\n  "node_address": "localhost:8080",\n  "address": "localhost:50052",\n  "listen_addr": ":50052",\n  "timeout": "5m",\n  "security": {\n    "mode": "none",\n    "cert_dir": "/etc/serviceradar/certs",\n    "role": "checker"\n  }\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"network-sweep",children:"Network Sweep"}),"\n",(0,i.jsxs)(n.p,{children:["For network scanning, edit ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/checkers/sweep/sweep.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "networks": ["192.168.2.0/24", "192.168.3.1/32"],\n  "ports": [22, 80, 443, 3306, 5432, 6379, 8080, 8443],\n  "sweep_modes": ["icmp", "tcp"],\n  "interval": "5m",\n  "concurrency": 100,\n  "timeout": "10s"\n}\n'})}),"\n",(0,i.jsx)(n.h2,{id:"next-steps",children:"Next Steps"}),"\n",(0,i.jsx)(n.p,{children:"After configuring your components:"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsx)(n.li,{children:"Restart services to apply changes:"}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"sudo systemctl restart serviceradar-agent\nsudo systemctl restart serviceradar-poller\nsudo systemctl restart serviceradar-core\n"})}),"\n",(0,i.jsxs)(n.ol,{start:"2",children:["\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsxs)(n.p,{children:["Visit the web dashboard at ",(0,i.jsx)(n.code,{children:"http://core-host:8090"})]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsxs)(n.p,{children:["Review ",(0,i.jsx)(n.a,{href:"/docs/tls-security",children:"TLS Security"})," to secure your components"]}),"\n"]}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,o.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(a,{...e})}):a(e)}},3881:(e,n,r)=>{r.d(n,{R:()=>c,x:()=>t});var s=r(8101);const i={},o=s.createContext(i);function c(e){const n=s.useContext(o);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function t(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:c(e.components),s.createElement(o.Provider,{value:n},e.children)}}}]);