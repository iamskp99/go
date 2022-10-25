package main

import (
	"context"
	"es4gophers/domain"
	"es4gophers/logic"
	"fmt"
)

func main() {

	ctx := context.Background()

	ctx = logic.LoadMoviesFromFile(ctx)
	fmt.Println(ctx.Value(domain.MoviesKey))
	// ctx = logic.ConnectWithElasticsearch(ctx)
	// logic.IndexMoviesAsDocuments(ctx)
	// logic.QueryMovieByDocumentID(ctx)
	// logic.BestKeanuActionMovies(ctx)
	// logic.MovieCountPerGenreAgg(ctx)

}
