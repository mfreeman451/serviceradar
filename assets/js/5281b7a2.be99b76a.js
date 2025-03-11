"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[2443],{3881:(e,r,n)=>{n.d(r,{R:()=>l,x:()=>o});var i=n(8101);const t={},s=i.createContext(t);function l(e){const r=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(r):{...r,...e}}),[r,e])}function o(e){let r;return r=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:l(e.components),i.createElement(s.Provider,{value:r},e.children)}},4798:(e,r,n)=>{n.r(r),n.d(r,{assets:()=>c,contentTitle:()=>o,default:()=>h,frontMatter:()=>l,metadata:()=>i,toc:()=>d});const i=JSON.parse('{"id":"architecture","title":"Architecture","description":"ServiceRadar uses a distributed, multi-layered architecture designed for flexibility, reliability, and security. This page explains how the different components work together to provide robust monitoring capabilities.","source":"@site/docs/architecture.md","sourceDirName":".","slug":"/architecture","permalink":"/docs/architecture","draft":false,"unlisted":false,"tags":[],"version":"current","sidebarPosition":6,"frontMatter":{"sidebar_position":6,"title":"Architecture"},"sidebar":"tutorialSidebar","previous":{"title":"Web UI Configuration","permalink":"/docs/web-ui"}}');var t=n(5105),s=n(3881);const l={sidebar_position:6,title:"Architecture"},o="ServiceRadar Architecture",c={},d=[{value:"Architecture Overview",id:"architecture-overview",level:2},{value:"Key Components",id:"key-components",level:2},{value:"Agent (Monitored Host)",id:"agent-monitored-host",level:3},{value:"Poller (Monitoring Coordinator)",id:"poller-monitoring-coordinator",level:3},{value:"Core Service (API &amp; Processing)",id:"core-service-api--processing",level:3},{value:"Web UI (User Interface)",id:"web-ui-user-interface",level:3},{value:"Security Architecture",id:"security-architecture",level:2},{value:"mTLS Security",id:"mtls-security",level:3},{value:"API Authentication",id:"api-authentication",level:3},{value:"Deployment Models",id:"deployment-models",level:2},{value:"Standard Deployment",id:"standard-deployment",level:3},{value:"Minimal Deployment",id:"minimal-deployment",level:3},{value:"High Availability Deployment",id:"high-availability-deployment",level:3},{value:"Network Requirements",id:"network-requirements",level:2}];function a(e){const r={a:"a",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",mermaid:"mermaid",p:"p",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,s.R)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(r.header,{children:(0,t.jsx)(r.h1,{id:"serviceradar-architecture",children:"ServiceRadar Architecture"})}),"\n",(0,t.jsx)(r.p,{children:"ServiceRadar uses a distributed, multi-layered architecture designed for flexibility, reliability, and security. This page explains how the different components work together to provide robust monitoring capabilities."}),"\n",(0,t.jsx)(r.h2,{id:"architecture-overview",children:"Architecture Overview"}),"\n",(0,t.jsx)(r.mermaid,{value:'graph TD\n    subgraph "User Access"\n        Browser[Web Browser]\n    end\n\n    subgraph "Service Layer"\n        WebUI[Web UI<br>:80/nginx]\n        CoreAPI[Core Service<br>:8090/:50052]\n        WebUI --\x3e|API calls<br>w/key auth| CoreAPI\n        Browser --\x3e|HTTP/HTTPS| WebUI\n    end\n\n    subgraph "Monitoring Layer"\n        Poller1[Poller 1<br>:50053]\n        Poller2[Poller 2<br>:50053]\n        CoreAPI ---|gRPC<br>bidirectional| Poller1\n        CoreAPI ---|gRPC<br>bidirectional| Poller2\n    end\n\n    subgraph "Target Infrastructure"\n        Agent1[Agent 1<br>:50051]\n        Agent2[Agent 2<br>:50051]\n        Agent3[Agent 3<br>:50051]\n        \n        Poller1 ---|gRPC<br>checks| Agent1\n        Poller1 ---|gRPC<br>checks| Agent2\n        Poller2 ---|gRPC<br>checks| Agent3\n        \n        Agent1 --- Service1[Services<br>Processes<br>Ports]\n        Agent2 --- Service2[Services<br>Processes<br>Ports]\n        Agent3 --- Service3[Services<br>Processes<br>Ports]\n    end\n\n    subgraph "Alerting"\n        CoreAPI --\x3e|Webhooks| Discord[Discord]\n        CoreAPI --\x3e|Webhooks| Other[Other<br>Services]\n    end\n\n    style Browser fill:#f9f,stroke:#333,stroke-width:1px\n    style WebUI fill:#b9c,stroke:#333,stroke-width:1px\n    style CoreAPI fill:#9bc,stroke:#333,stroke-width:2px\n    style Poller1 fill:#adb,stroke:#333,stroke-width:1px\n    style Poller2 fill:#adb,stroke:#333,stroke-width:1px\n    style Agent1 fill:#fd9,stroke:#333,stroke-width:1px\n    style Agent2 fill:#fd9,stroke:#333,stroke-width:1px\n    style Agent3 fill:#fd9,stroke:#333,stroke-width:1px\n    style Discord fill:#c9d,stroke:#333,stroke-width:1px\n    style Other fill:#c9d,stroke:#333,stroke-width:1px'}),"\n",(0,t.jsx)(r.h2,{id:"key-components",children:"Key Components"}),"\n",(0,t.jsx)(r.h3,{id:"agent-monitored-host",children:"Agent (Monitored Host)"}),"\n",(0,t.jsx)(r.p,{children:"The Agent runs on each host you want to monitor and is responsible for:"}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Collecting service status information (process status, port availability, etc.)"}),"\n",(0,t.jsx)(r.li,{children:"Exposing a gRPC service on port 50051 for Pollers to query"}),"\n",(0,t.jsx)(r.li,{children:"Supporting various checker types (process, port, SNMP, etc.)"}),"\n",(0,t.jsx)(r.li,{children:"Running with minimal privileges for security"}),"\n"]}),"\n",(0,t.jsx)(r.p,{children:(0,t.jsx)(r.strong,{children:"Technical Details:"})}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Written in Go for performance and minimal dependencies"}),"\n",(0,t.jsx)(r.li,{children:"Uses gRPC for efficient, language-agnostic communication"}),"\n",(0,t.jsx)(r.li,{children:"Supports dynamic loading of checker plugins"}),"\n",(0,t.jsx)(r.li,{children:"Can run on constrained hardware with minimal resource usage"}),"\n"]}),"\n",(0,t.jsx)(r.h3,{id:"poller-monitoring-coordinator",children:"Poller (Monitoring Coordinator)"}),"\n",(0,t.jsx)(r.p,{children:"The Poller coordinates monitoring activities and is responsible for:"}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Querying multiple Agents at configurable intervals"}),"\n",(0,t.jsx)(r.li,{children:"Aggregating status data from Agents"}),"\n",(0,t.jsx)(r.li,{children:"Reporting status to the Core Service"}),"\n",(0,t.jsx)(r.li,{children:"Performing direct checks (HTTP, ICMP, etc.)"}),"\n",(0,t.jsx)(r.li,{children:"Supporting network sweeps and discovery"}),"\n"]}),"\n",(0,t.jsx)(r.p,{children:(0,t.jsx)(r.strong,{children:"Technical Details:"})}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Runs on port 50053 for gRPC communications"}),"\n",(0,t.jsx)(r.li,{children:"Stateless design allows multiple Pollers for high availability"}),"\n",(0,t.jsx)(r.li,{children:"Configurable polling intervals for different check types"}),"\n",(0,t.jsx)(r.li,{children:"Supports both pull-based (query) and push-based (events) monitoring"}),"\n"]}),"\n",(0,t.jsx)(r.h3,{id:"core-service-api--processing",children:"Core Service (API & Processing)"}),"\n",(0,t.jsx)(r.p,{children:"The Core Service is the central component that:"}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Receives and processes reports from Pollers"}),"\n",(0,t.jsx)(r.li,{children:"Provides an API for the Web UI on port 8090"}),"\n",(0,t.jsx)(r.li,{children:"Triggers alerts based on configurable thresholds"}),"\n",(0,t.jsx)(r.li,{children:"Stores historical monitoring data"}),"\n",(0,t.jsx)(r.li,{children:"Manages webhook notifications"}),"\n"]}),"\n",(0,t.jsx)(r.p,{children:(0,t.jsx)(r.strong,{children:"Technical Details:"})}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Exposes a gRPC service on port 50052 for Poller connections"}),"\n",(0,t.jsx)(r.li,{children:"Provides a RESTful API on port 8090 for the Web UI"}),"\n",(0,t.jsx)(r.li,{children:"Uses role-based security model"}),"\n",(0,t.jsx)(r.li,{children:"Implements webhook templating for flexible notifications"}),"\n"]}),"\n",(0,t.jsx)(r.h3,{id:"web-ui-user-interface",children:"Web UI (User Interface)"}),"\n",(0,t.jsx)(r.p,{children:"The Web UI provides a modern dashboard interface that:"}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Visualizes the status of monitored services"}),"\n",(0,t.jsx)(r.li,{children:"Displays historical performance data"}),"\n",(0,t.jsx)(r.li,{children:"Provides configuration management"}),"\n",(0,t.jsx)(r.li,{children:"Securely communicates with the Core Service API"}),"\n"]}),"\n",(0,t.jsx)(r.p,{children:(0,t.jsx)(r.strong,{children:"Technical Details:"})}),"\n",(0,t.jsxs)(r.ul,{children:["\n",(0,t.jsx)(r.li,{children:"Built with Next.js in SSR mode for security and performance"}),"\n",(0,t.jsx)(r.li,{children:"Uses Nginx as a reverse proxy on port 80"}),"\n",(0,t.jsx)(r.li,{children:"Communicates with the Core Service API using a secure API key"}),"\n",(0,t.jsx)(r.li,{children:"Supports responsive design for mobile and desktop"}),"\n"]}),"\n",(0,t.jsx)(r.h2,{id:"security-architecture",children:"Security Architecture"}),"\n",(0,t.jsx)(r.p,{children:"ServiceRadar implements multiple layers of security:"}),"\n",(0,t.jsx)(r.h3,{id:"mtls-security",children:"mTLS Security"}),"\n",(0,t.jsx)(r.p,{children:"For network communication between components, ServiceRadar supports mutual TLS (mTLS):"}),"\n",(0,t.jsx)(r.mermaid,{value:'graph TB\nsubgraph "Agent Node"\nAG[Agent<br/>Role: Server<br/>:50051]\nSNMPCheck[SNMP Checker<br/>:50054]\nDuskCheck[Dusk Checker<br/>:50052]\nSweepCheck[Network Sweep]\n\n        AG --\x3e SNMPCheck\n        AG --\x3e DuskCheck\n        AG --\x3e SweepCheck\n    end\n    \n    subgraph "Poller Service"\n        PL[Poller<br/>Role: Client+Server<br/>:50053]\n    end\n    \n    subgraph "Core Service"\n        CL[Core Service<br/>Role: Server<br/>:50052]\n        DB[(Database)]\n        API[HTTP API<br/>:8090]\n        \n        CL --\x3e DB\n        CL --\x3e API\n    end\n    \n    %% Client connections from Poller\n    PL --\x3e|mTLS Client| AG\n    PL --\x3e|mTLS Client| CL\n    \n    %% Server connections to Poller\n    HC1[Health Checks] --\x3e|mTLS Client| PL\n    \n    classDef server fill:#e1f5fe,stroke:#01579b\n    classDef client fill:#f3e5f5,stroke:#4a148c\n    classDef dual fill:#fff3e0,stroke:#e65100\n    \n    class AG,CL server\n    class PL dual\n    class SNMPCheck,DuskCheck,SweepCheck client'}),"\n",(0,t.jsx)(r.h3,{id:"api-authentication",children:"API Authentication"}),"\n",(0,t.jsx)(r.p,{children:"The Web UI communicates with the Core Service using API key authentication:"}),"\n",(0,t.jsx)(r.mermaid,{value:"sequenceDiagram\n    participant User as User (Browser)\n    participant WebUI as Web UI (Next.js)\n    participant API as Core API\n    \n    User->>WebUI: HTTP Request\n    Note over WebUI: Server-side middleware<br>loads API key\n    WebUI->>API: Request with API Key\n    API->>API: Validate API Key\n    API->>WebUI: Response\n    WebUI->>User: Rendered UI"}),"\n",(0,t.jsxs)(r.p,{children:["For details on configuring security, see the ",(0,t.jsx)(r.a,{href:"/docs/tls-security",children:"TLS Security"})," documentation."]}),"\n",(0,t.jsx)(r.h2,{id:"deployment-models",children:"Deployment Models"}),"\n",(0,t.jsx)(r.p,{children:"ServiceRadar supports multiple deployment models:"}),"\n",(0,t.jsx)(r.h3,{id:"standard-deployment",children:"Standard Deployment"}),"\n",(0,t.jsx)(r.p,{children:"All components installed on separate machines for optimal security and reliability:"}),"\n",(0,t.jsx)(r.mermaid,{value:"graph LR\n    Browser[Browser] --\x3e WebServer[Web Server<br/>Web UI + Core]\n    WebServer --\x3e PollerServer[Poller Server]\n    PollerServer --\x3e AgentServer1[Host 1<br/>Agent]\n    PollerServer --\x3e AgentServer2[Host 2<br/>Agent]\n    PollerServer --\x3e AgentServerN[Host N<br/>Agent]"}),"\n",(0,t.jsx)(r.h3,{id:"minimal-deployment",children:"Minimal Deployment"}),"\n",(0,t.jsx)(r.p,{children:"For smaller environments, components can be co-located:"}),"\n",(0,t.jsx)(r.mermaid,{value:"graph LR\n    Browser[Browser] --\x3e CombinedServer[Combined Server<br/>Web UI + Core + Poller]\n    CombinedServer --\x3e AgentServer1[Host 1<br/>Agent]\n    CombinedServer --\x3e AgentServer2[Host 2<br/>Agent]"}),"\n",(0,t.jsx)(r.h3,{id:"high-availability-deployment",children:"High Availability Deployment"}),"\n",(0,t.jsx)(r.p,{children:"For mission-critical environments:"}),"\n",(0,t.jsx)(r.mermaid,{value:"graph TD\n    LB[Load Balancer] --\x3e WebServer1[Web Server 1<br/>Web UI]\n    LB --\x3e WebServer2[Web Server 2<br/>Web UI]\n    WebServer1 --\x3e CoreServer1[Core Server 1]\n    WebServer2 --\x3e CoreServer1\n    WebServer1 --\x3e CoreServer2[Core Server 2]\n    WebServer2 --\x3e CoreServer2\n    CoreServer1 --\x3e Poller1[Poller 1]\n    CoreServer2 --\x3e Poller1\n    CoreServer1 --\x3e Poller2[Poller 2]\n    CoreServer2 --\x3e Poller2\n    Poller1 --\x3e Agent1[Agent 1]\n    Poller1 --\x3e Agent2[Agent 2]\n    Poller2 --\x3e Agent1\n    Poller2 --\x3e Agent2"}),"\n",(0,t.jsx)(r.h2,{id:"network-requirements",children:"Network Requirements"}),"\n",(0,t.jsx)(r.p,{children:"ServiceRadar uses the following network ports:"}),"\n",(0,t.jsxs)(r.table,{children:[(0,t.jsx)(r.thead,{children:(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.th,{children:"Component"}),(0,t.jsx)(r.th,{children:"Port"}),(0,t.jsx)(r.th,{children:"Protocol"}),(0,t.jsx)(r.th,{children:"Purpose"})]})}),(0,t.jsxs)(r.tbody,{children:[(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Agent"}),(0,t.jsx)(r.td,{children:"50051"}),(0,t.jsx)(r.td,{children:"gRPC/TCP"}),(0,t.jsx)(r.td,{children:"Service status queries"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Poller"}),(0,t.jsx)(r.td,{children:"50053"}),(0,t.jsx)(r.td,{children:"gRPC/TCP"}),(0,t.jsx)(r.td,{children:"Health checks"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Core"}),(0,t.jsx)(r.td,{children:"50052"}),(0,t.jsx)(r.td,{children:"gRPC/TCP"}),(0,t.jsx)(r.td,{children:"Poller connections"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Core"}),(0,t.jsx)(r.td,{children:"8090"}),(0,t.jsx)(r.td,{children:"HTTP/TCP"}),(0,t.jsx)(r.td,{children:"API (internal)"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Web UI"}),(0,t.jsx)(r.td,{children:"80/443"}),(0,t.jsx)(r.td,{children:"HTTP(S)/TCP"}),(0,t.jsx)(r.td,{children:"User interface"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"SNMP Checker"}),(0,t.jsx)(r.td,{children:"50054"}),(0,t.jsx)(r.td,{children:"gRPC/TCP"}),(0,t.jsx)(r.td,{children:"SNMP status queries"})]}),(0,t.jsxs)(r.tr,{children:[(0,t.jsx)(r.td,{children:"Dusk Checker"}),(0,t.jsx)(r.td,{children:"50052"}),(0,t.jsx)(r.td,{children:"gRPC/TCP"}),(0,t.jsx)(r.td,{children:"Dusk node monitoring"})]})]})]}),"\n",(0,t.jsxs)(r.p,{children:["For more information on deploying ServiceRadar, see the ",(0,t.jsx)(r.a,{href:"/docs/installation",children:"Installation Guide"}),"."]})]})}function h(e={}){const{wrapper:r}={...(0,s.R)(),...e.components};return r?(0,t.jsx)(r,{...e,children:(0,t.jsx)(a,{...e})}):a(e)}}}]);