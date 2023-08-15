import React from 'react';
import { useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import useMediaQuery from '@mui/material/useMediaQuery';
import Main from 'layouts/Main';

const NotFoundCover = (): JSX.Element => {
  const theme = useTheme();
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
            404
          </Typography>
          <Typography
            variant="h6"
            component="p"
            color="text.secondary"
            align={isMd ? 'left' : 'center'}
          >
            Oops! Looks like you followed a bad link.
            <br />
            If you think this is a problem with us, please{' '}
            <Link href={''} underline="none">
              tell us
            </Link>
          </Typography>
          <Box
            marginTop={4}
            display={'flex'}
            justifyContent={{ xs: 'center', md: 'flex-start' }}
          >
            <Button
              component={Link}
              variant="contained"
              color="primary"
              size="large"
              href={'/'}
            >
              Back home
            </Button>
          </Box>
        </Box>
      </Box>
    </Main>
  );
};

export default NotFoundCover;
