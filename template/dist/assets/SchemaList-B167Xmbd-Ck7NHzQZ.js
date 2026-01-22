import{k as g,S as N,l as y,a as e,_ as c,B as o,m as v,J as m,o as l,C as b,p as w,s as d,u as S,w as C,x as z,y as $}from"./entry.client-Cy9IsQOv.js";import{z as p}from"./index-DRBOFufT-C5jqVKPr.js";import{R as A}from"./Toc-PbuF-u9x-ijQYm_0d.js";import{F as P,x as R}from"./SchemaView-Dx9gtne2-Be40Zbai.js";import"./circular-B5aq06Lb-BwXoodZ8.js";const I=$(`
  query GetSchemas($input: JSON!, $type: SchemaType!) {
    schema(input: $input, type: $type) {
      title
      description
      summary
      components {
        schemas {
          name
          schema
          extensions
        }
      }
    }
  }
`);function D(){const{input:h,type:x,versions:u,version:t,options:r}=g(),j=N(I,{input:h,type:x}),{data:a}=y(j),n=a.schema.components?.schemas??[],f=Object.entries(u).length>1,i=r?.showVersionSelect==="always"||f&&r?.showVersionSelect!=="hide";return n.length?e.jsxs("div",{className:"grid grid-cols-(--sidecar-grid-cols) gap-8 justify-between","data-pagefind-filter":"section:openapi","data-pagefind-meta":"section:openapi",children:[e.jsx(P,{name:"category",children:a.schema.title}),e.jsxs(c,{children:[e.jsxs("title",{children:["Schemas ",i?t:""]}),e.jsx("meta",{name:"description",content:"List of schemas used by the API."})]}),e.jsxs("div",{className:"pt-(--padding-content-top) pb-(--padding-content-bottom)",children:[e.jsxs(o,{className:"w-full",children:[e.jsxs("div",{className:"flex flex-col gap-y-4 sm:flex-row justify-around items-start sm:items-end",children:[e.jsxs("div",{className:"flex-1",children:[e.jsx(v,{children:a.schema.title}),e.jsxs(m,{level:1,id:"schemas",registerNavigationAnchor:!0,className:"mb-0",children:["Schemas",i&&e.jsxs("span",{className:"text-xl text-muted-foreground ms-1.5",children:["(",t,")"]})]})]}),a.schema.description&&e.jsxs(l,{className:"flex items-center gap-1 text-sm font-medium text-muted-foreground group",children:[e.jsx("span",{children:"API information"}),e.jsx(b,{className:"group-data-[state=open]:hidden translate-y-px",size:14}),e.jsx(w,{className:"group-data-[state=closed]:hidden translate-y-px",size:13})]})]}),a.schema.description&&e.jsx(d,{className:"CollapsibleContent",children:e.jsx("div",{className:"mt-4 max-w-full border rounded-sm bg-muted/25",children:e.jsx(S,{className:"max-w-full prose-img:max-w-prose border-border p-3 lg:p-5",content:a.schema.description})})})]}),e.jsx("hr",{className:"my-8"}),e.jsx("div",{className:"flex flex-col gap-y-5",children:n.map(s=>e.jsxs(o,{className:"group",defaultOpen:!0,children:[e.jsxs(m,{registerNavigationAnchor:!0,level:2,className:"flex items-center gap-1 justify-between w-fit",id:p(s.name),children:[s.name," ",e.jsx(l,{asChild:!0,children:e.jsx(C,{variant:"ghost",size:"icon",className:"size-6",children:e.jsx(z,{size:16,className:"group-data-[state=open]:rotate-90 transition cursor-pointer"})})})]}),e.jsx(d,{className:"mt-4 CollapsibleContent",children:e.jsx(R,{schema:s.schema})})]},s.name))})]}),e.jsx(A,{entries:n.map(s=>({id:p(s.name),value:s.name,depth:1}))})]}):e.jsxs("div",{children:[e.jsxs(c,{children:[e.jsxs("title",{children:["Schemas ",i?t:""]}),e.jsx("meta",{name:"description",content:"List of schemas used by the API."})]}),"No schemas found"]})}export{D as SchemaList};
//# sourceMappingURL=SchemaList-B167Xmbd-Ck7NHzQZ.js.map
