"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[2887],{609:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>a,contentTitle:()=>t,default:()=>h,frontMatter:()=>l,metadata:()=>s,toc:()=>o});const s=JSON.parse('{"id":"tls-security","title":"TLS Security","description":"ServiceRadar supports mutual TLS (mTLS) authentication to secure communications between components. This guide explains how to set up and configure mTLS for your ServiceRadar deployment.","source":"@site/docs/tls-security.md","sourceDirName":".","slug":"/tls-security","permalink":"/serviceradar/docs/tls-security","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":4,"frontMatter":{"sidebar_position":4,"title":"TLS Security"},"sidebar":"tutorialSidebar","previous":{"title":"Configuration Basics","permalink":"/serviceradar/docs/configuration"},"next":{"title":"Web UI Configuration","permalink":"/serviceradar/docs/web-ui"}}');var i=r(5105),c=r(3881);const l={sidebar_position:4,title:"TLS Security"},t="TLS Security",a={},o=[{value:"Security Architecture",id:"security-architecture",level:2},{value:"Certificate Overview",id:"certificate-overview",level:2},{value:"Certificate Generation",id:"certificate-generation",level:2},{value:"1. Install cfssl",id:"1-install-cfssl",level:3},{value:"2. Create Configuration Files",id:"2-create-configuration-files",level:3},{value:"3. Generate Certificates",id:"3-generate-certificates",level:3},{value:"Certificate Deployment",id:"certificate-deployment",level:2},{value:"Role-Based Requirements",id:"role-based-requirements",level:3},{value:"Installation Steps",id:"installation-steps",level:3},{value:"Expected Directory Structure",id:"expected-directory-structure",level:3},{value:"Component Configuration",id:"component-configuration",level:2},{value:"Agent Configuration",id:"agent-configuration",level:3},{value:"Poller Configuration",id:"poller-configuration",level:3},{value:"Core Configuration",id:"core-configuration",level:3},{value:"Verification",id:"verification",level:2},{value:"Troubleshooting",id:"troubleshooting",level:2}];function d(e){const n={admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",mermaid:"mermaid",ol:"ol",p:"p",pre:"pre",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,c.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.header,{children:(0,i.jsx)(n.h1,{id:"tls-security",children:"TLS Security"})}),"\n",(0,i.jsx)(n.p,{children:"ServiceRadar supports mutual TLS (mTLS) authentication to secure communications between components. This guide explains how to set up and configure mTLS for your ServiceRadar deployment."}),"\n",(0,i.jsx)(n.h2,{id:"security-architecture",children:"Security Architecture"}),"\n",(0,i.jsx)(n.p,{children:"ServiceRadar components communicate securely using mTLS with the following roles:"}),"\n",(0,i.jsx)(n.mermaid,{value:'graph TB\nsubgraph "Agent Node"\nAG[Agent<br/>Role: Server<br/>:50051]\nSNMPCheck[SNMP Checker<br/>:50054]\nDuskCheck[Dusk Checker<br/>:50052]\nSweepCheck[Network Sweep]\n\n        AG --\x3e SNMPCheck\n        AG --\x3e DuskCheck\n        AG --\x3e SweepCheck\n    end\n    \n    subgraph "Poller Service"\n        PL[Poller<br/>Role: Client+Server<br/>:50053]\n    end\n    \n    subgraph "Core Service"\n        CL[Core Service<br/>Role: Server<br/>:50052]\n        DB[(Database)]\n        API[HTTP API<br/>:8090]\n        \n        CL --\x3e DB\n        CL --\x3e API\n    end\n    \n    %% Client connections from Poller\n    PL --\x3e|mTLS Client| AG\n    PL --\x3e|mTLS Client| CL\n    \n    %% Server connections to Poller\n    HC1[Health Checks] --\x3e|mTLS Client| PL\n    \n    classDef server fill:#e1f5fe,stroke:#01579b\n    classDef client fill:#f3e5f5,stroke:#4a148c\n    classDef dual fill:#fff3e0,stroke:#e65100\n    \n    class AG,CL server\n    class PL dual\n    class SNMPCheck,DuskCheck,SweepCheck client'}),"\n",(0,i.jsx)(n.h2,{id:"certificate-overview",children:"Certificate Overview"}),"\n",(0,i.jsx)(n.p,{children:"ServiceRadar uses the following certificate files:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"root.pem"})," - The root CA certificate"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"server.pem"})," - Server certificate"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"server-key.pem"})," - Server private key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"client.pem"})," - Client certificate"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"client-key.pem"})," - Client private key"]}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"certificate-generation",children:"Certificate Generation"}),"\n",(0,i.jsx)(n.h3,{id:"1-install-cfssl",children:"1. Install cfssl"}),"\n",(0,i.jsx)(n.p,{children:"First, install the CloudFlare SSL toolkit (cfssl) which we'll use for generating certificates:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"go install github.com/cloudflare/cfssl/cmd/...@latest\n"})}),"\n",(0,i.jsx)(n.h3,{id:"2-create-configuration-files",children:"2. Create Configuration Files"}),"\n",(0,i.jsxs)(n.p,{children:["Create ",(0,i.jsx)(n.code,{children:"cfssl.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n    "signing": {\n        "default": {\n            "expiry": "8760h"\n        },\n        "profiles": {\n            "rootca": {\n                "usages": ["signing", "key encipherment", "server auth", "client auth"],\n                "expiry": "8760h",\n                "ca_constraint": {\n                    "is_ca": true,\n                    "max_path_len": 0\n                }\n            },\n            "server": {\n                "usages": ["signing", "key encipherment", "server auth"],\n                "expiry": "8760h"\n            },\n            "client": {\n                "usages": ["signing", "key encipherment", "client auth"],\n                "expiry": "8760h"\n            }\n        }\n    }\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["Create ",(0,i.jsx)(n.code,{children:"csr.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n    "hosts": ["localhost", "127.0.0.1"],\n    "key": {\n        "algo": "ecdsa",\n        "size": 256\n    },\n    "names": [{\n        "O": "ServiceRadar"\n    }]\n}\n'})}),"\n",(0,i.jsx)(n.admonition,{type:"note",children:(0,i.jsx)(n.p,{children:'Modify the "hosts" array in csr.json to include the actual hostnames and IP addresses of your ServiceRadar components.'})}),"\n",(0,i.jsx)(n.h3,{id:"3-generate-certificates",children:"3. Generate Certificates"}),"\n",(0,i.jsx)(n.p,{children:"Generate the root CA:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:'cfssl selfsign -config cfssl.json --profile rootca "ServiceRadar CA" csr.json | cfssljson -bare root\n'})}),"\n",(0,i.jsx)(n.p,{children:"Generate server and client keys:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"cfssl genkey csr.json | cfssljson -bare server\ncfssl genkey csr.json | cfssljson -bare client\n"})}),"\n",(0,i.jsx)(n.p,{children:"Sign the certificates:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile server server.csr | cfssljson -bare server\ncfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile client client.csr | cfssljson -bare client\n"})}),"\n",(0,i.jsx)(n.h2,{id:"certificate-deployment",children:"Certificate Deployment"}),"\n",(0,i.jsx)(n.h3,{id:"role-based-requirements",children:"Role-Based Requirements"}),"\n",(0,i.jsx)(n.p,{children:"Different ServiceRadar components need different certificates based on their role:"}),"\n",(0,i.jsxs)(n.table,{children:[(0,i.jsx)(n.thead,{children:(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.th,{children:"Component"}),(0,i.jsx)(n.th,{children:"Role"}),(0,i.jsx)(n.th,{children:"Certificates Needed"})]})}),(0,i.jsxs)(n.tbody,{children:[(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:"Poller"}),(0,i.jsx)(n.td,{children:"Client+Server"}),(0,i.jsx)(n.td,{children:"All certificates (client + server)"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:"Agent"}),(0,i.jsx)(n.td,{children:"Client+Server"}),(0,i.jsx)(n.td,{children:"All certificates (client + server)"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:"Core"}),(0,i.jsx)(n.td,{children:"Server only"}),(0,i.jsx)(n.td,{children:"Server certificates only"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:"Checker"}),(0,i.jsx)(n.td,{children:"Server only"}),(0,i.jsx)(n.td,{children:"Server certificates only"})]})]})]}),"\n",(0,i.jsx)(n.h3,{id:"installation-steps",children:"Installation Steps"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsx)(n.li,{children:"Create the certificates directory:"}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"sudo mkdir -p /etc/serviceradar/certs\nsudo chown serviceradar:serviceradar /etc/serviceradar/certs\nsudo chmod 700 /etc/serviceradar/certs\n"})}),"\n",(0,i.jsxs)(n.ol,{start:"2",children:["\n",(0,i.jsx)(n.li,{children:"Install certificates based on role:"}),"\n"]}),"\n",(0,i.jsx)(n.p,{children:"For core/checker (server-only):"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"sudo cp root.pem server*.pem /etc/serviceradar/certs/\n"})}),"\n",(0,i.jsx)(n.p,{children:"For poller/agent (full set):"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"sudo cp root.pem server*.pem client*.pem /etc/serviceradar/certs/\n"})}),"\n",(0,i.jsxs)(n.ol,{start:"3",children:["\n",(0,i.jsx)(n.li,{children:"Set permissions:"}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"sudo chown serviceradar:serviceradar /etc/serviceradar/certs/*\nsudo chmod 644 /etc/serviceradar/certs/*.pem\nsudo chmod 600 /etc/serviceradar/certs/*-key.pem\n"})}),"\n",(0,i.jsx)(n.h3,{id:"expected-directory-structure",children:"Expected Directory Structure"}),"\n",(0,i.jsx)(n.p,{children:"Server-only deployment (Core/Checker):"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"/etc/serviceradar/certs/\n\u251c\u2500\u2500 root.pem           (644)\n\u251c\u2500\u2500 server.pem         (644)\n\u2514\u2500\u2500 server-key.pem     (600)\n"})}),"\n",(0,i.jsx)(n.p,{children:"Full deployment (Poller/Agent):"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"/etc/serviceradar/certs/\n\u251c\u2500\u2500 root.pem           (644)\n\u251c\u2500\u2500 server.pem         (644)\n\u251c\u2500\u2500 server-key.pem     (600)\n\u251c\u2500\u2500 client.pem         (644)\n\u2514\u2500\u2500 client-key.pem     (600)\n"})}),"\n",(0,i.jsx)(n.h2,{id:"component-configuration",children:"Component Configuration"}),"\n",(0,i.jsx)(n.h3,{id:"agent-configuration",children:"Agent Configuration"}),"\n",(0,i.jsxs)(n.p,{children:["Update ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/agent.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "checkers_dir": "/etc/serviceradar/checkers",\n  "listen_addr": ":50051",\n  "service_type": "grpc",\n  "service_name": "AgentService",\n  "security": {\n    "mode": "mtls",\n    "cert_dir": "/etc/serviceradar/certs",\n    "server_name": "poller-hostname",\n    "role": "agent"\n  }\n}\n'})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Set ",(0,i.jsx)(n.code,{children:"mode"})," to ",(0,i.jsx)(n.code,{children:"mtls"})]}),"\n",(0,i.jsxs)(n.li,{children:["Set ",(0,i.jsx)(n.code,{children:"server_name"})," to the hostname/IP address of the poller"]}),"\n",(0,i.jsxs)(n.li,{children:["Set ",(0,i.jsx)(n.code,{children:"role"})," to ",(0,i.jsx)(n.code,{children:"agent"})]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"poller-configuration",children:"Poller Configuration"}),"\n",(0,i.jsxs)(n.p,{children:["Update ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/poller.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "agents": {\n    "local-agent": {\n      "address": "agent-hostname:50051",\n      "security": {\n        "server_name": "agent-hostname",\n        "mode": "mtls"\n      },\n      "checks": [\n        // your checks here\n      ]\n    }\n  },\n  "core_address": "core-hostname:50052",\n  "listen_addr": ":50053",\n  "poll_interval": "30s",\n  "poller_id": "my-poller",\n  "service_name": "PollerService",\n  "service_type": "grpc",\n  "security": {\n    "mode": "mtls",\n    "cert_dir": "/etc/serviceradar/certs",\n    "server_name": "core-hostname",\n    "role": "poller"\n  }\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"core-configuration",children:"Core Configuration"}),"\n",(0,i.jsxs)(n.p,{children:["Update ",(0,i.jsx)(n.code,{children:"/etc/serviceradar/core.json"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "listen_addr": ":8090",\n  "grpc_addr": ":50052",\n  "alert_threshold": "5m",\n  "known_pollers": ["my-poller"],\n  "metrics": {\n    "enabled": true,\n    "retention": 100,\n    "max_nodes": 10000\n  },\n  "security": {\n    "mode": "mtls",\n    "cert_dir": "/etc/serviceradar/certs",\n    "role": "core"\n  }\n}\n'})}),"\n",(0,i.jsx)(n.h2,{id:"verification",children:"Verification"}),"\n",(0,i.jsx)(n.p,{children:"Verify your installation with:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"ls -la /etc/serviceradar/certs/\n"})}),"\n",(0,i.jsx)(n.p,{children:"Example output (Core instance):"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{children:"total 20\ndrwx------ 2 serviceradar serviceradar 4096 Feb 21 20:46 .\ndrwxr-xr-x 3 serviceradar serviceradar 4096 Feb 21 22:35 ..\n-rw-r--r-- 1 serviceradar serviceradar  920 Feb 21 04:47 root.pem\n-rw------- 1 serviceradar serviceradar  227 Feb 21 20:44 server-key.pem\n-rw-r--r-- 1 serviceradar serviceradar  928 Feb 21 20:45 server.pem\n"})}),"\n",(0,i.jsx)(n.h2,{id:"troubleshooting",children:"Troubleshooting"}),"\n",(0,i.jsx)(n.p,{children:"Common issues and solutions:"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Permission denied"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Verify directory permissions: ",(0,i.jsx)(n.code,{children:"700"})]}),"\n",(0,i.jsxs)(n.li,{children:["Verify file permissions: ",(0,i.jsx)(n.code,{children:"644"})," for certificates, ",(0,i.jsx)(n.code,{children:"600"})," for keys"]}),"\n",(0,i.jsxs)(n.li,{children:["Check ownership: ",(0,i.jsx)(n.code,{children:"serviceradar:serviceradar"})]}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Certificate not found"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"Confirm all required certificates for the role are present"}),"\n",(0,i.jsx)(n.li,{children:"Double-check file paths in configuration"}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Invalid certificate"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"Ensure certificates are properly formatted PEM files"}),"\n",(0,i.jsx)(n.li,{children:"Verify certificates were generated with correct profiles"}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Connection refused"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"Verify server name in config matches certificate CN"}),"\n",(0,i.jsx)(n.li,{children:"Check that all required certificates are present for the role"}),"\n",(0,i.jsx)(n.li,{children:"Confirm service has proper read permissions"}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsx)(n.p,{children:(0,i.jsx)(n.strong,{children:"Testing with grpcurl"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Install grpcurl: ",(0,i.jsx)(n.code,{children:"go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"})]}),"\n",(0,i.jsxs)(n.li,{children:["Test health check endpoint with mTLS:","\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"grpcurl -cacert root.pem \\\n        -cert client.pem \\\n        -key client-key.pem \\\n        -servername <SERVER_IP> \\\n        <SERVER_IP>:50052 \\\n        grpc.health.v1.Health/Check\n"})}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["Successful response should show:","\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "status": "SERVING"\n}\n'})}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["If this fails but certificates are correct, verify:","\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"Server is running and listening on specified port"}),"\n",(0,i.jsx)(n.li,{children:"Firewall rules allow the connection"}),"\n",(0,i.jsx)(n.li,{children:"The servername matches the certificate's Common Name (CN)"}),"\n"]}),"\n"]}),"\n"]}),"\n"]}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,c.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},3881:(e,n,r)=>{r.d(n,{R:()=>l,x:()=>t});var s=r(8101);const i={},c=s.createContext(i);function l(e){const n=s.useContext(c);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function t(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:l(e.components),s.createElement(c.Provider,{value:n},e.children)}}}]);