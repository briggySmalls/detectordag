import { shallow } from "enzyme";
import React from "react";
import renderer from "react-test-renderer";
import axios from 'axios';

import Home from "../pages/index.js";
import Layout from "../components/layout.js";
import Form from '../components/registrationForm';

// Mock axios for us
jest.mock('axios');

describe("With Enzyme", () => {
  it('Home has expected content', () => {
    const app = shallow(<Home />);
    expect(app.find("Layout")).toHaveLength(1);
    expect(app.find("Layout Alert")).toHaveLength(0);
    expect(app.find("Layout p").text()).toEqual("Register your device to get started");
  });

  it('Home calls API successfully', () => {
    const app = shallow(<Home />);
    // Simulate a submission
    const testData = {test: 'data'};
    app.find(Form).props().onSubmit(testData);
    // Expect axios to have been called
    const resp = {data: users};
    axios.post.mockResolvedValue(resp);
    // Assert
    expect(axios.post.mock.calls.length).toBe(1);
    expect(axios.post.mock.calls[0][0]).toBe('/api/register');
    expect(axios.post.mock.calls[0][1]).toBe(testData);
  });

  it('Home calls API which fails', () => {
    const app = shallow(<Home />);
    // Simulate a submission
    const testData = {test: 'data'};
    app.find(Form).props().onSubmit(testData);
    // Expect axios to have been called
    axios.post.mockRejectedValue(new Error("A failure"));
    // Assert mock calls
    expect(axios.post.mock.calls.length).toBe(1);
    expect(axios.post.mock.calls[0][0]).toBe('/api/register');
    expect(axios.post.mock.calls[0][1]).toBe(testData);
    // Assert error
    const error = app.find("Error");
    expect(error).toHaveLength(1);
    expect(error.text()).toBe("A failure");
    expect(error).prop('severity').toEqual('error');
  });
});

describe("With Snapshot Testing", () => {
  it('Home shows "Hello, Sunshine!"', () => {
    const component = renderer.create(<Home />);
    const tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
