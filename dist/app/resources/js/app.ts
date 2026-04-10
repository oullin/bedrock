import { createApp, h, DefineComponent } from "vue";
import { createInertiaApp } from "@inertiajs/vue3";

import "../css/app.css";

createInertiaApp({
  title: (title) => (title ? `${title} - Bedrock` : "Bedrock"),
  resolve: (name) => {
    const pages = import.meta.glob<DefineComponent>("./Pages/**/*.vue", { eager: true });
    return pages[`./Pages/${name}.vue`];
  },
  setup({ el, App, props, plugin }) {
    createApp({ render: () => h(App, props) })
      .use(plugin)
      .mount(el);
  },
});
