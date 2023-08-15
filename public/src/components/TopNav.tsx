import React from 'react';
import Box from '@mui/material/Box';

import ThemeModeToggler from 'components/ThemeModeToggler';

const TopNav = (): JSX.Element => {
  return (
    <Box display={'flex'} justifyContent={'flex-end'} alignItems={'center'}>
      <Box>
        <ThemeModeToggler />
      </Box>
    </Box>
  );
};

export default TopNav;
