import { primary_theme } from '../../util/Theme';
import { ThemeProvider, StyledEngineProvider } from '@mui/material';
import { withRouter } from 'storybook-addon-react-router-v6';
import { Provider } from 'react-redux';
import { store } from '../../Store';
import RecentActivities from '../../Components/RecentActivitiesList/RecentActivities';

export default {
  title: 'ProfilePage/RecentActivities',
  component: RecentActivities,
  decorators: [
    withRouter,
    (Story) => (
      <Provider store={store}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={primary_theme}>
            <Story />
          </ThemeProvider>
        </StyledEngineProvider>
      </Provider>
    ),
  ],
  // This component will have an automatically generated Autodocs entry: https://storybook.js.org/docs/writing-docs/autodocs
  tags: ['autodocs'],
  parameters: {
    // More on how to position stories at: https://storybook.js.org/docs/configure/story-layout
    layout: 'fullscreen',
  },
};

export const PrimaryRecentActivities = {
  args: {},
};
