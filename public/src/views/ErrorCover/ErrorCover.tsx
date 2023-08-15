import React from 'react';
import { useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import useMediaQuery from '@mui/material/useMediaQuery';
import Main from 'layouts/Main';
import { useNavigate } from 'react-router-dom';

const ErrorCover = (): JSX.Element => {
  const theme = useTheme();
  const navigate = useNavigate();
  const isMd = useMediaQuery(theme.breakpoints.up('md'), {
    defaultMatches: true
  });

  return (
    <Main>
      <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
        <Box marginTop={4} marginBottom={4}>
          <Typography
            variant="h1"
            component={'h1'}
            align={isMd ? 'left' : 'center'}
            sx={{ fontWeight: 700 }}
          >
            Oops!
          </Typography>
          <Typography
            variant="h6"
            component="p"
            color="text.secondary"
            align={isMd ? 'left' : 'center'}
          >
            Looks like something went wrong.
            <br />
            This seems to be on our side, give it some seconds and try again
            please.
          </Typography>
          <Box
            marginTop={4}
            display={'flex'}
            justifyContent={{ xs: 'center', md: 'flex-start' }}
          >
            <Button
              onClick={() => navigate(-1)}
              variant="outlined"
              color="primary"
              size="large"
              sx={{ marginRight: 2 }}
            >
              Go back
            </Button>
            <Button
              component={Link}
              variant="contained"
              color="primary"
              size="large"
              href={'/'}
            >
              Go home
            </Button>
          </Box>
        </Box>
      </Box>
    </Main>
  );
};

export default ErrorCover;
