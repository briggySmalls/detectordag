import { shallow } from "enzyme";
import React from "react";
import renderer from "react-test-renderer";
import axios from 'axios';

import Home from "../pages/index.jsx";
import Layout from "../components/layout.jsx";
import Form from '../components/registrationForm';

const mockPush = jest.fn().mockName('router.push');
jest.mock("next/router", () => ({useRouter() {return {push: mockPush}}}));

// Mock axios for us
jest.mock('axios');

describe("With Enzyme", () => {
  beforeEach(() => {
    // Reset mocks
    axios.post.mockReset();
    mockPush.mockReset();
  });

  it('Home has expected content', () => {
    const app = shallow(<Home />);
    expect(app.find("Layout")).toHaveLength(1);
    expect(app.find("Layout Alert")).toHaveLength(0);
    expect(app.find("Layout p").text()).toEqual("Register your device to get started");
  });

  it('Home calls API successfully', async () => {
    const app = shallow(<Home />);
    const loading = app.find("WithLoadingComponent");
    // Ensure the loading value is false
    expect(loading.prop('isLoading')).toBe(false);
    // Configure axios mock success
    axios.post.mockImplementationOnce(() => {
        // Assert that the loading icon is shown
        // expect(loading.props().isLoading).toBe(true);
        // Return response
        return Promise.resolve({data: 'whatever really'});
    });
    // Simulate a submission
    const testData = {test: 'data'};
    loading.prop('onSubmit')(testData);
    // Assert
    await expect(axios.post).toHaveBeenCalledWith('/api/register', testData);
    expect(loading.prop('isLoading')).toBe(false);
    console.log(loading.debug());
    expect(mockPush).toHaveBeenCalledWith('/success');
  });

  it('Home calls API which fails', async () => {
    const app = shallow(<Home />);
    const loading = app.find("WithLoadingComponent");
    // Configure axios mock failure
    axios.post.mockRejectedValueOnce(new Error("A failure"));
    // Simulate a submission
    const testData = {test: 'data'};
    loading.prop('onSubmit')(testData);
    // Assert mock calls
    await expect(axios.post).toHaveBeenCalledWith('/api/register', testData);
    expect(mockPush).not.toHaveBeenCalled();
    // Assert error
    const error = app.find("Alert");
    expect(error).toHaveLength(1);
    expect(error.text()).toBe("A failure");
    expect(error.prop('severity')).toEqual('error');
  });
});

describe("With Snapshot Testing", () => {
  it('Home shows "Hello, Sunshine!"', () => {
    const component = renderer.create(<Home />);
    const tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
