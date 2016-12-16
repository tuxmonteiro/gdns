require 'test_helper'

class Bind9ControllerTest < ActionDispatch::IntegrationTest
  test "should get scheduler_export" do
    get bind9_scheduler_export_url
    assert_response :success
  end

  test "should get export" do
    get bind9_export_url
    assert_response :success
  end

end
