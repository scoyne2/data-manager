import Grid from "@mui/material/Grid";
import FormNewFeed from "src/views/form-layouts/FormNewFeed";

const FormLayouts = () => {
  return (
    <Grid container spacing={12}>
      <Grid item xs={12} md={12}>
        <p></p>
        <FormNewFeed />
      </Grid>
    </Grid>
  );
};

export default FormLayouts;
