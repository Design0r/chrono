import daisyui from "daisyui";
export default {
  content: ["../**/*.html", "../**/*.templ", "../**/*.go"],
  theme: {},
  plugins: [daisyui],
  daisyui: {
    themes: [
      {
        dim: {
          ...require("daisyui/src/theming/themes")["dim"],
          info: "#9ED0F1",
          error: "#fe7c5d",
        },
      },
    ],
  },
};
