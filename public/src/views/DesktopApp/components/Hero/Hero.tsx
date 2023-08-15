import React from 'react';
import { useTheme } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';

const Hero = (): JSX.Element => {
  const theme = useTheme();

  const isMd = useMediaQuery(theme.breakpoints.up('md'), {
    defaultMatches: true
  });

  return (
    <Grid container spacing={4}>
      <Grid item container xs={12} md={6} alignItems={'center'}>
        <Box data-aos={isMd ? 'fade-right' : 'fade-up'}>
          <Box marginBottom={2}>
            <Typography
              variant="h3"
              color="text.primary"
              sx={{ fontWeight: 700 }}
            >
              Generating{' '}
              <Link href={'https://www.apple.com/wallet/'} target={'_blank'}>
                Wallet Passes
              </Link>{' '}
              has never been so easy
            </Typography>
          </Box>
          <Box marginBottom={3}>
            <Typography variant="h6" component="p" color="text.secondary">
              Use ChatGPT to generate your own Wallet Passes. No coding
              required. Generate one or multiple passes at once, with barcodes,
              images, and more. Install the plugin and say "Generate a Wallet
              Pass" to get started.
            </Typography>
          </Box>
        </Box>
      </Grid>
      <Grid
        item
        container
        alignItems={'center'}
        justifyContent={'center'}
        xs={12}
        md={6}
        data-aos="flip-left"
        data-aos-easing="ease-out-cubic"
        data-aos-duration="2000"
      >
        <Box
          component={'img'}
          loading="lazy"
          src={'https://storage.googleapis.com/walletpasses/usage2.png'}
          alt="Use ChatGPT to generate Wallet Passes"
          boxShadow={3}
          borderRadius={2}
          sx={{
            filter: theme.palette.mode === 'dark' ? 'brightness(0.7)' : 'none'
          }}
        />
      </Grid>
    </Grid>
  );
};

export default Hero;
