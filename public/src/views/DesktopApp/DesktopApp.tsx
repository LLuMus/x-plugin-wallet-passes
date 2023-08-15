import React from 'react';
import { useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Main from 'layouts/Main';
import Container from 'components/Container';
import { Hero } from './components';
import Prompts from './components/Prompts';
import Typography from '@mui/material/Typography';

const DesktopApp = (): JSX.Element => {
  const theme = useTheme();
  return (
    <Main>
      <Box
        position={'relative'}
        sx={{
          backgroundColor: theme.palette.alternate.main,
          marginTop: -13,
          paddingTop: 13
        }}
      >
        <Container>
          <Hero />
        </Container>
        <Box marginTop={10}>
          <Box marginBottom={2}>
            <Typography
              variant="h4"
              color="text.primary"
              align={'center'}
              gutterBottom
              sx={{
                fontWeight: 700
              }}
            >
              Here are some prompts to get you started
            </Typography>
            <Typography
              variant="h6"
              component="p"
              color="text.secondary"
              sx={{ fontWeight: 400 }}
              align={'center'}
            >
              You can click on the Copy Prompt button to copy the prompt and
              then paste on ChatGPT.
            </Typography>
          </Box>
        </Box>
        <Container>
          <Prompts />
        </Container>
      </Box>
    </Main>
  );
};

export default DesktopApp;
