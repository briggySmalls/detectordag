// Mock axios for us
jest.mock('axios');
const mockPush = jest.fn();
jest.mock("next/router", () => ({
    useRouter() {
        return {
            push: mockPush,
        };
    },
}));

import { shallow } from "enzyme";
import React from "react";
import renderer from "react-test-renderer";
import axios from 'axios';

import Home from "../pages/index.jsx";
import Layout from "../components/layout.jsx";
import Form from '../components/registrationForm';

describe("With Enzyme", () => {
  beforeEach(() => {
    axios.post.mockReset();
  });

  it('Home has expected content', () => {
    const app = shallow(<Home />);
    expect(app.find("Layout")).toHaveLength(1);
    expect(app.find("Layout Alert")).toHaveLength(0);
    expect(app.find("Layout p").text()).toEqual("Register your device to get started");
  });

  it('Home calls API successfully', () => {
    const app = shallow(<Home />);
    // Configure axios mock success
    const resp = {data: 'whatever really'};
    axios.post.mockResolvedValue(resp);
    // Simulate a submission
    const testData = {test: 'data'};
    app.find("WithLoadingComponent").props().onSubmit(testData);
    // Assert
    expect(axios.post).toHaveBeenCalledWith('/api/register', testData)
    // expect(mockPush).toHaveBeenCalledWith('/success')
  });

  it('Home calls API which fails', () => {
    const app = shallow(<Home />);
    // Configure axios mock failure
    axios.post.mockRejectedValue(new Error("A failure"));
    // Simulate a submission
    const testData = {test: 'data'};
    app.find("WithLoadingComponent").props().onSubmit(testData);
    // Assert mock calls
    expect(axios.post).toHaveBeenCalledWith('/api/register', testData)
    // Assert error
    console.log(app.debug());
    const error = app.find("Alert");
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
