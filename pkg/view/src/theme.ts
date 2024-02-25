import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
  palette: {
    primary: {
      main: "#49509e",
    },
    secondary: {
      main: "#d43630",
    },
  },
  typography: {
    h1: { fontSize: "2.5rem" },
    h2: { fontSize: "2rem" },
    h3: { fontSize: "1.75rem" },
    h4: { fontSize: "1.5rem" },
    h5: { fontSize: "1.25rem" },
    h6: { fontSize: "1.15rem" },
  },
});
