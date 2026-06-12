package converter

import (
	"ultrathreads/model"
	"ultrathreads/util/hashid"
)

func ToTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	slug := hashid.Id2Slug[model.Tag](tag.ID)
	return &model.TagResponse{Slug: slug, TagName: tag.Name}
}

func ToTags(tags []model.Tag) []model.TagResponse {
	if len(tags) == 0 {
		return []model.TagResponse{}
	}
	responses := make([]model.TagResponse, 0, len(tags))
	for i := range tags {
		slug := hashid.Id2Slug[model.Tag](tags[i].ID)
		responses = append(responses, model.TagResponse{
			Slug:    slug,
			TagName: tags[i].Name,
		})
	}
	return responses
}