package converter

import (
	"ultrathreads/model"
)

func ToTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	return &model.TagResponse{TagId: tag.ID, TagName: tag.Name}
}

func ToTags(tags []model.Tag) []model.TagResponse {
	if len(tags) == 0 {
		return []model.TagResponse{}
	}
	responses := make([]model.TagResponse, 0, len(tags))
	for i := range tags {
		responses = append(responses, model.TagResponse{
			TagId:   tags[i].ID,
			TagName: tags[i].Name,
		})
	}
	return responses
}