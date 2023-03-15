import StateReport from './components/StateReport'
import './App.css';
import VerticalSlider from './components/VerticalSlider';
import { Typography, AppBar, Button, Card, CardActions, CardContent, CardMedia, CssBaseline, Grid, Toolbar, Container } from '@mui/material';
import { Castle, PrecisionManufacturing } from '@mui/icons-material';

const pieces = [
  { num: 1 },
  { num: 2 },
  { num: 3 },
  { num: 4 },
]

function App() {
  return (
    <>
      <CssBaseline />
      <AppBar position="relative">
        <Toolbar>
          <PrecisionManufacturing/>
          <Typography variant="h6">
            A Study of Light
          </Typography>
        </Toolbar>
      </AppBar>
      <Container maxWidth="sm">
        <Typography
          variant="h2"
          align="center"
          color="textPrimary"
          gutterBottom
        >
            Local Control Interface
        </Typography>
        <div>
          <Grid container spacing={2} justifyContent="center">
            <Grid item>
              <Button variant="contained" color="primary">Wake</Button>
            </Grid>
            <Grid item>
              <Button variant="contained" color="secondary">Shutdown</Button>
            </Grid>
          </Grid>
        </div>
      </Container>
      <Container>
        <Grid container spacing={4}>
          {pieces.map((piece) => (
            <Grid item>
              <Card>
                <CardMedia sx={{height:"200px"}} image="https://source.unsplash.com/random" title="title"/>
                <CardContent>
                  <Typography variant="h4" align="center">{piece.num}</Typography>
                </CardContent>
                <CardActions>
                  <Button size="small" color="primary">View</Button>
                  <Button size="small" color="secondary">Edit</Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>
      <StateReport/>
      <VerticalSlider onValueChange={console.log}/>
    </>
  );
}

export default App;
