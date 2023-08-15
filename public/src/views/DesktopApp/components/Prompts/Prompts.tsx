import React from 'react';
import { useTheme } from '@mui/material/styles';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Grid from '@mui/material/Grid';
import { Card, CardContent, Chip } from '@mui/material';
import Button from '@mui/material/Button';
import { ContentCopy } from '@mui/icons-material';

const mock = [
  {
    image: 'https://storage.googleapis.com/walletpasses/sample1.png',
    description: 'Please generate a sample Wallet Pass use this website for sample images https://picsum.photos/200/200',
    title: 'Your first pass!',
    tags: ['Beginner', 'Example', 'Tutorial']
  },
  {
    image: 'https://storage.googleapis.com/walletpasses/sample2.png',
    description:
      'Could you please help me generate a Wallet Pass for my private health insurance? ' +
      'Here are the personal details that should be included: - Name: John Doe - Insurance ' +
      'Number: XXX5521215523 (please include this as a secondary field) - Logo: https://t.ly/8C3nd ' +
      '- Logo Text: HanseMerkur - Height: 1.85cm - Born: 1999 - Blood Type: A+ - Member Since: 2017 - ' +
      'Background Color: Use this green #027040. Please note that a barcode should not be included in the ' +
      'pass. On the back of the Wallet Pass, could you also include the link to the HanseMerkur website ' +
      '(https://www.hansemerkur.de/) and the telephone number of my consultant (+49 999 9999)? Thank you!',
    title: 'Insurance Details',
    tags: ['Easy', 'Insurance', 'Utility']
  },
  {
    image: 'https://storage.googleapis.com/walletpasses/sample3.png',
    description:
      'Could you please help me generate a coupon with a discount for my Wallmart branch? Here are the details:' +
      ' - Logo: https://t.ly/xU2ls - Logo Text: Wallmart - Strip: https://t.ly/-3Ada - Send an empty primary field,' +
      ' and secondary field the expiration date (2 weeks) - Here is the code for the barcode: SecretSecret123123 -' +
      ' Use this background color: #1A75CF. Thank you!',
    title: 'Customer Coupon',
    tags: ['Advanced', 'Quick', 'Coupon']
  }
];

const Prompts = (): JSX.Element => {
  const theme = useTheme();

  return (
    <Grid container spacing={4}>
      {mock.map((item, i) => (
        <Grid key={i} item xs={12}>
          <Box
            component={Card}
            width={1}
            height={1}
            borderRadius={0}
            boxShadow={0}
            display={'flex'}
            flexDirection={{
              xs: 'column',
              md: i % 2 === 0 ? 'row-reverse' : 'row'
            }}
            sx={{ backgroundImage: 'none', bgcolor: 'transparent' }}
          >
            <Box
              sx={{
                width: { xs: 1, md: '50%' }
              }}
            >
              <Box
                component={'img'}
                height={1}
                width={1}
                src={item.image}
                alt="..."
                sx={{
                  objectFit: 'contain',
                  maxHeight: 620,
                  filter:
                    theme.palette.mode === 'dark' ? 'brightness(0.7)' : 'none'
                }}
              />
            </Box>
            <CardContent
              sx={{
                paddingX: { xs: 1, sm: 2, md: 4 },
                paddingY: { xs: 2, sm: 4 },
                width: { xs: 1, md: '50%' },
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center',
                alignItems: i % 2 === 0 ? 'end' : 'start'
              }}
            >
              <Box>
                {item.tags.map(item => (
                  <Chip
                    key={item}
                    label={item}
                    component="a"
                    href=""
                    clickable
                    size={'small'}
                    color={'primary'}
                    sx={{ marginBottom: 1, marginRight: 1 }}
                  />
                ))}
              </Box>
              <Typography
                variant={'h6'}
                fontWeight={700}
                sx={{ textTransform: 'uppercase' }}
              >
                {item.title}
              </Typography>
              <Typography color="text.secondary">{item.description}</Typography>
              <Box marginTop={2} display={'flex'} justifyContent={'flex-end'}>
                <Button
                  endIcon={<ContentCopy />}
                  onClick={() => {
                    void navigator.clipboard.writeText(item.description);
                  }}
                >
                  Copy Prompt
                </Button>
              </Box>
            </CardContent>
          </Box>
        </Grid>
      ))}
    </Grid>
  );
};

export default Prompts;
