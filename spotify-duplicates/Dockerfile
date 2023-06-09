FROM mcr.microsoft.com/dotnet/runtime:7.0 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["spotify-duplicates/spotify-duplicates.csproj", "spotify-duplicates/"]
RUN dotnet restore "spotify-duplicates/spotify-duplicates.csproj"
COPY . .
WORKDIR "/src/spotify-duplicates"
RUN dotnet build "spotify-duplicates.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "spotify-duplicates.csproj" -c Release -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "spotify-duplicates.dll"]
