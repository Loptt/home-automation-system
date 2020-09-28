import React from "react";
import HomeScreen from "./src/views/HomeView";
import { createAppContainer, createSwitchNavigator } from "react-navigation";
import { createDrawerNavigator } from "react-navigation-drawer";
import { createStackNavigator } from "react-navigation-stack";
import { Dimensions } from "react-native";

import { Main } from "./views/index";

import MainView from "./views/MainView";

const AuthStack = createStackNavigator(
  {
    AuthView,
    LoginView,
    RegisterView,
  },
  {
    initialRouteName: "AuthView",
    headerMode: "none",
    navigationOptions: {
      headerShown: false,
    },
  }
);

const ItemStack = createStackNavigator(
  {
    ItemView,
  },
  {
    initialRouteName: "ItemView",
    headerMode: "none",
    navigationOptions: {
      headerShown: false,
    },
  }
);
/**
 * DrawerNavigator
 * Basically this class is used to define the navigator and the app container.
 * @author Luis F. Alvarez Sanchez
 * 8/29/2020
 */
const DrawerNavigator: any = createDrawerNavigator(
  {
    Main: {
      screen: MainView,
      navigationOptions: {
        title: "Home",
      },
    },
  },
  {
    contentComponent: (props) => <SideBar {...props} />,
    drawerWidth: Dimensions.get("window").width * 0.85,
    hideStatusBar: true,
  }
);

export default createAppContainer(
  createSwitchNavigator(
    {
      App: DrawerNavigator,
      Auth: AuthStack,
      Item: ItemStack,
    },
    {
      initialRouteName: "Auth",
    }
  )
);