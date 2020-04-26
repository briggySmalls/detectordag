import { shallow } from "enzyme";
import React from "react";
import renderer from "react-test-renderer";

import Home from "../pages/index.js";
import Layout from "../components/layout.js";

describe("With Enzyme", () => {
  it('Home shows "Hello, Sunshine!"', () => {
    const app = shallow(<Home />);
    expect(app.find("Layout")).toHaveLength(1);
    expect(app.find("Layout p").text()).toEqual("Register your device to get started");
  });
});

describe("With Snapshot Testing", () => {
  it('Home shows "Hello, Sunshine!"', () => {
    const component = renderer.create(<Home />);
    const tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
