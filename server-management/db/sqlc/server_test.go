package db

import (
	"context"
	"testing"
	"vcssms/util"

	"github.com/stretchr/testify/require"
)

func createRandomServer(t *testing.T) Server {
	server := CreateServerParams{
		Name:   util.RandomServerName(),
		Status: int32(util.RandomStatus()),
		Ipv4:   util.RandomIP(),
	}
	s, err := testQueries.CreateServer(context.Background(), server)
	require.NoError(t, err)
	require.NotZero(t, s)

	require.Equal(t, server.Name, s.Name)
	require.Equal(t, server.Ipv4, s.Ipv4)
	require.Equal(t, server.Status, s.Status)

	return s
}

func TestCreateServer(t *testing.T) {
	createRandomServer(t)
}

func TestGetServer(t *testing.T) {
	newServer := createRandomServer(t)
	server := GetServerParams{
		Limit:  1,
		Offset: 0,
		ID:     int32(newServer.ID),
		IsID:   true,
		IDAsc:  true,
	}
	s, err := testQueries.GetServer(context.Background(), server)
	require.NoError(t, err)
	require.NotEmpty(t, s)

	require.Equal(t, newServer.Name, s[0].Name)
	require.Equal(t, newServer.Ipv4, s[0].Ipv4)
	require.Equal(t, newServer.Status, s[0].Status)
}

func TestUpdateServer(t *testing.T) {
	newServer := createRandomServer(t)

	agr := UpdateServerParams{
		ID:      newServer.ID,
		Name:    "VCSabc",
		SetName: true,
	}

	s, err := testQueries.UpdateServer(context.Background(), agr)

	require.NoError(t, err)
	require.NotEmpty(t, s)

	require.Equal(t, s.Name, "VCSabc")
}

func TestDeleteAccount(t *testing.T) {
	newServer := createRandomServer(t)

	err := testQueries.DeleteServer(context.Background(), newServer.ID)
	require.NoError(t, err)

	s, err := testQueries.GetServer(context.Background(), GetServerParams{
		Limit:  1,
		Offset: 0,
		ID:     int32(newServer.ID),
		IsID:   true,
	})
	require.Empty(t, s)
}
