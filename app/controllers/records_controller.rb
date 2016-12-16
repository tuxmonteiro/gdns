class RecordsController < ApplicationController
  before_action :set_record, only: [:show, :update, :destroy]

  # GET /domains/1/records
  def index
    begin
    @records = Record.where(domain_id: params[:domain_id])
      render json: @records
    rescue ActiveRecord::RecordNotFound
      render :head => true, :status => :not_found
    end
  end

  # GET /domains/1/records/1
  def show
    if @record.nil?
      render :head => true, :status => :not_found
    else
      render json: @record
    end
  end

  # POST /records
  def create
    @record = Record.new(record_params)

    if @record.save
      render json: @record, status: :created, location: @record
    else
      render json: @record.errors, status: :unprocessable_entity
    end
  end

  # PATCH/PUT /records/1
  def update
    if @record.update(record_params)
      render json: @record
    else
      render json: @record.errors, status: :unprocessable_entity
    end
  end

  # DELETE /records/1
  def destroy
    @record.destroy
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_record
      begin
        @record = Record.where(domain_id: params[:domain_id]).where(id: params[:id]).first!
      rescue ActiveRecord::RecordNotFound
        @record = nil
      end
    end

    # Only allow a trusted parameter "white list" through.
    def record_params
      params.require(:record).permit(:name, :type, :content, :domain_id)
    end
end
