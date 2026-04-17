import{d as b,j as i,w as s,g as t,a as o,u as r,o as d,n as f,b as a,f as g}from"./app.js";import{_,a as n}from"./FeatureHeader.vue_vue_type_script_setup_true_lang-4CZFCx4c.js";import{f as m,_ as c}from"./routes-S6IyX_wi.js";import{_ as l}from"./index-BdaDjCGN.js";import{_ as v}from"./AppLayout.vue_vue_type_script_setup_true_lang-Ct7coe5q.js";import"./CardTitle.vue_vue_type_script_setup_true_lang-ucxQeZZT.js";import"./CardDescription.vue_vue_type_script_setup_true_lang-dHSpRoiL.js";const h={class:"flex flex-1 flex-col gap-6 p-4 lg:p-6"},w={class:"grid gap-6 lg:grid-cols-2"},y={class:"space-y-4"},P={class:"flex flex-wrap gap-2"},k={class:"space-y-3 text-sm"},C={class:"space-y-2"},B={class:"flex items-center gap-2 rounded-md border p-2"},S={class:"flex items-center gap-2 rounded-md border p-2"},I={class:"flex items-center gap-2 rounded-md border p-2"},E=b({__name:"Progress",setup(N){const x=[{title:"Features"},{title:"Events"},{title:"Progress"}],p=m("features.data-loading.deferred-props"),u=m("features.data-loading.polling");return(D,e)=>(d(),i(v,{title:"Progress",breadcrumbs:x},{default:s(()=>[t("div",h,[o(_,{title:"Progress Indicator",description:"Inertia shows a progress bar during page navigations. The loading bar at the top of this app demonstrates this feature."}),t("div",w,[o(n,{title:"Try It",description:"Navigate to see the progress bar in action."},{default:s(()=>[t("div",y,[e[2]||(e[2]=t("p",{class:"text-sm"}," Click a link below and watch the loading bar at the top of the page. Slower endpoints make the progress bar more visible. ",-1)),t("div",P,[r(p)?(d(),i(r(c),{key:0,"as-child":"",variant:"outline"},{default:s(()=>[o(r(f),{href:r(p)},{default:s(()=>[...e[0]||(e[0]=[a(" Deferred Props (slow) ",-1)])]),_:1},8,["href"])]),_:1})):g("",!0),r(u)?(d(),i(r(c),{key:1,"as-child":"",variant:"outline"},{default:s(()=>[o(r(f),{href:r(u)},{default:s(()=>[...e[1]||(e[1]=[a(" Polling Page ",-1)])]),_:1},8,["href"])]),_:1})):g("",!0)])])]),_:1}),o(n,{title:"Configuration",description:"How the progress bar is configured."},{default:s(()=>[...e[3]||(e[3]=[t("div",{class:"space-y-3 text-sm"},[t("pre",{class:"bg-muted overflow-auto rounded p-4 text-xs leading-relaxed"},`createInertiaApp({
  progress: {
    // Delay before showing (ms)
    delay: 250,

    // Color of the bar
    color: '#4B5563',

    // Show the spinner
    includeCSS: true,

    // Show spinner icon
    showSpinner: false,
  },
})`),t("p",{class:"text-muted-foreground"}," The progress bar only appears after a short delay to avoid flickering on fast navigations. ")],-1)])]),_:1}),o(n,{title:"Custom Progress Bar",description:"This app uses a custom LoadingBar component."},{default:s(()=>[t("div",k,[e[10]||(e[10]=t("p",null,[a(" Instead of using the built-in NProgress bar, this demo app uses a custom "),t("code",{class:"bg-muted rounded px-1.5 py-0.5 text-xs"},"LoadingBar"),a(" component that listens to router events. ")],-1)),t("div",C,[t("div",B,[o(r(l),{variant:"secondary"},{default:s(()=>[...e[4]||(e[4]=[a("start",-1)])]),_:1}),e[5]||(e[5]=t("span",{class:"text-muted-foreground text-xs"},"Show the progress bar",-1))]),t("div",S,[o(r(l),{variant:"secondary"},{default:s(()=>[...e[6]||(e[6]=[a("progress",-1)])]),_:1}),e[7]||(e[7]=t("span",{class:"text-muted-foreground text-xs"},"Update the bar width",-1))]),t("div",I,[o(r(l),{variant:"secondary"},{default:s(()=>[...e[8]||(e[8]=[a("finish",-1)])]),_:1}),e[9]||(e[9]=t("span",{class:"text-muted-foreground text-xs"},"Complete and hide the bar",-1))])])])]),_:1}),o(n,{title:"Disabling Progress",description:"Opt out of the progress indicator."},{default:s(()=>[...e[11]||(e[11]=[t("div",{class:"space-y-3 text-sm"},[t("p",null,"You can disable the progress bar globally or per-visit:"),t("pre",{class:"bg-muted overflow-auto rounded p-3 text-xs leading-relaxed"},`// Disable globally
createInertiaApp({
  progress: false,
})

// Disable per-visit
router.visit(url, {
  showProgress: false,
})`)],-1)])]),_:1})])])]),_:1}))}});export{E as default};
